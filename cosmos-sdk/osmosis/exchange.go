package osmosis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	codecType "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdksigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	osmo "github.com/osmosis-labs/osmosis/v11/x/gamm/types"
	tmjson "github.com/tendermint/tendermint/libs/json"

	cosmossdk "github.com/rootwarp/cosmev/cosmos-sdk"
	"github.com/rootwarp/cosmev/types"
)

type osmosisExchanger struct {
	kr      keyring.Keyring
	keyInfo keyring.Info

	rpcURL string

	authClient      cosmossdk.Auth
	accountNo       uint64
	currentSequence uint64
}

func (e *osmosisExchanger) SetMnemonic(mnemonic, hdPath string) error {
	e.kr = keyring.NewInMemory()

	sigAlgo, err := keyring.NewSigningAlgoFromString("secp256k1", keyring.SigningAlgoList{hd.Secp256k1})
	if err != nil {
		return err
	}

	e.keyInfo, err = e.kr.NewAccount("swap", mnemonic, keyring.DefaultBIP39Passphrase, hdPath, sigAlgo)
	if err != nil {
		return err
	}

	e.authClient = NewAuthClient(e.rpcURL)
	account, err := e.authClient.GetAccount(e.keyInfo.GetAddress().String())
	if err != nil {
		return err
	}

	accountNo, err := strconv.ParseInt(account.AccountNumber, 10, 64)
	if err != nil {
		return err
	}

	e.accountNo = uint64(accountNo)

	seq, err := strconv.ParseInt(account.Sequence, 10, 64)
	if err != nil {
		return err
	}

	e.currentSequence = uint64(seq)

	log.Println("Account is ready", account.Address, account.AccountNumber, account.Sequence)

	return nil
}

func (e *osmosisExchanger) Address() string {
	return e.keyInfo.GetAddress().String()
}

func (e *osmosisExchanger) Swap(tokenIn types.Asset, routes []types.Pool, minTokenOutAmount int64) (*types.Asset, error) {
	in := sdk.Coin{
		Amount: sdk.NewIntFromBigInt(tokenIn.Amount),
		Denom:  tokenIn.Denom,
	}

	// TODO:
	// #1: osmo -> atom, fee 0.002
	// #498: atom -> juno fee 0.003
	// #497: juno -> osmo fee  0.003
	wrapRoutes := []osmo.SwapAmountInRoute{
		{
			PoolId:        1,
			TokenOutDenom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", // uatom
		},
		{
			PoolId:        498,
			TokenOutDenom: "ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED", // ujuno
		},
		{
			PoolId:        497,
			TokenOutDenom: "uosmo",
		},
	}

	// TODO:
	fees := sdk.Coins{{Amount: sdk.NewInt(1250), Denom: "uosmo"}}
	txMsg, err := e.signTx(in, wrapRoutes, minTokenOutAmount, 250000, fees)

	fmt.Println("Final Tx Msg", string(txMsg))

	err = e.sendTx(txMsg)

	fmt.Println(err)

	return nil, err
}

// TODO:
func (e *osmosisExchanger) createTx(tokenIn sdk.Coin, routes []osmo.SwapAmountInRoute, minTokenOutAmount int64) (string, error) {
	return "", nil
}

func (e *osmosisExchanger) signTx(
	tokenIn sdk.Coin,
	routes []osmo.SwapAmountInRoute,
	minTokenOutAmount int64,
	gasLimit uint64,
	fees sdk.Coins,
) ([]byte, error) {

	ifRegistry := codecType.NewInterfaceRegistry()
	osmo.RegisterInterfaces(ifRegistry)

	marshaller := codec.NewProtoCodec(ifRegistry)
	txConfig := tx.NewTxConfig(marshaller, tx.DefaultSignModes)

	txBuilder := txConfig.NewTxBuilder()

	swapIn := osmo.MsgSwapExactAmountIn{
		Sender:            e.keyInfo.GetAddress().String(),
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: sdk.NewInt(minTokenOutAmount),
	}

	txBuilder.SetMsgs(&swapIn)
	txBuilder.SetMemo("")
	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetTimeoutHeight(0)
	txBuilder.SetFeeAmount(fees)

	jsonByte, err := txConfig.TxJSONEncoder()(txBuilder.GetTx())
	fmt.Println("JSON Msg ", string(jsonByte), err)

	signerData := signing.SignerData{
		ChainID:       "osmosis-1",
		AccountNumber: e.accountNo,
		Sequence:      e.currentSequence,
	}

	singleSigData := sdksigning.SingleSignatureData{
		SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}

	sigV2 := sdksigning.SignatureV2{
		PubKey:   e.keyInfo.GetPubKey(),
		Data:     &singleSigData,
		Sequence: e.currentSequence,
	}

	txBuilder.SetSignatures(sigV2)

	byteToSign, err := txConfig.
		SignModeHandler().
		GetSignBytes(sdksigning.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())

	// Signing here.
	sig, _, err := e.kr.Sign("swap", byteToSign)
	if err != nil {
		return nil, err
	}

	// Set Signature
	singleSigData = sdksigning.SingleSignatureData{
		SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
		Signature: sig,
	}

	sigV2 = sdksigning.SignatureV2{
		PubKey:   e.keyInfo.GetPubKey(),
		Data:     &singleSigData,
		Sequence: e.currentSequence,
	}

	txBuilder.SetSignatures(sigV2)
	//

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	rawTx, err := tmjson.Marshal(txBytes)
	if err != nil {
		return nil, err
	}

	rawMsg := map[string]json.RawMessage{}
	rawMsg["tx"] = rawTx

	rawParam, err := json.Marshal(rawMsg)
	if err != nil {
		return nil, err
	}

	rpcMsg := struct {
		JSONRPC string
		ID      int
		Method  string
		Params  json.RawMessage
	}{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "broadcast_tx_sync",
		Params:  rawParam,
	}

	msg, err := json.Marshal(rpcMsg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (e *osmosisExchanger) sendTx(txMsg []byte) error {
	cli := http.Client{}

	// TODO
	//url := "https://osmosis-mainnet-rpc.allthatnode.com:26657"
	//url := "https://rpc-osmosis.keplr.app"
	url := "https://osmosis-rpc.polkachu.com"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(txMsg))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

// TODO: Refactoring
type direction int

const (
	pathIn direction = iota
	pathOut
)

func (e *osmosisExchanger) convert(pools []types.Pool, inDenom string) ([]osmo.SwapAmountInRoute, error) {
	// TODO: Check len(pools)

	type swapflow struct {
		In  string
		Out string
	}

	flows := make([]swapflow, len(pools))

	// Start.

	first := pools[0]
	firstFlow := make([]direction, len(first.PoolAssets))

	curInDenom := inDenom
	for i, poolAsset := range first.PoolAssets {
		_, _, err := first.FindPoolAssetByDenom(inDenom)
		if err != nil {
			return nil, err
		}

		if poolAsset.Denom == inDenom {
			firstFlow[i] = pathIn
		} else {
			firstFlow[i] = pathOut
			curInDenom = poolAsset.Denom
		}
	}

	//poolAsset, err := first.FindPoolAssetByDenom(inDenom)
	//if err != nil {
	//	return nil, err
	//}

	//flows[0].In = poolAsset.Denom

	_ = flows
	_ = err

	return nil, nil
}

// NewExchanger returns DexExchanger instance.
func NewExchanger(rpc string) types.DexExchanger {
	return &osmosisExchanger{
		rpcURL: rpc,
	}
}
