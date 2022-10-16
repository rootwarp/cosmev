package osmosis

import (
	"math/big"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/rootwarp/cosmev/types"
	"github.com/stretchr/testify/assert"
)

/*
{
  "body": {
    "messages": [
      {
        "@type": "/osmosis.gamm.v1beta1.MsgSwapExactAmountIn",
        "sender": "osmo1nh6w4cm8scm8pvkmqhxm5m74pl73n987y3swaf",
        "routes": [
          {
            "poolId": "1",
            "tokenOutDenom": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"
          }
        ],
        "tokenIn": {
          "denom": "uosmo",
          "amount": "100000"
        },
        "tokenOutMinAmount": "15500"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [
        {
          "denom": "uosmo",
          "amount": "300"
        }
      ],
      "gas_limit": "250000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
*/

var (
	testMnemonic = os.Getenv("COSMEV_TEST_MNEMONIC")
	testAddress  = os.Getenv("COSMEV_TEST_ADDRESS")
)

// x/gamm/types/*
//func TestSwap(t *testing.T) {
//	kr := keyring.NewInMemory()
//	info, err := createKeys(kr, testMnemonic, "test")
//
//	assert.Nil(t, err)
//	assert.Equal(t, testAddress, info.GetAddress().String())
//
//	tokenIn := sdk.Coin{
//		Amount: sdk.NewInt(0),
//		Denom:  "uatom",
//	}
//
//	routes := []osmo.SwapAmountInRoute{
//		{
//			PoolId:        1,
//			TokenOutDenom: "uosmo",
//		},
//		{
//			PoolId:        2,
//			TokenOutDenom: "uion",
//		},
//	}
//
//	err = createMsg(kr, info, tokenIn, routes, 0, sdk.Coins{{Amount: sdk.NewInt(1), Denom: "uosmo"}})
//
//	fmt.Println(err)
//}

func TestExchange(t *testing.T) {
	exchanger := NewExchanger(os.Getenv("COSMEV_TEST_RPC"))

	hdParam := hd.CreateHDPath(118, 0, 0)
	err := exchanger.SetMnemonic(testMnemonic, hdParam.String())

	assert.Nil(t, err)
	assert.Equal(t, testAddress, exchanger.Address())

	in := types.Asset{
		Denom:  "uosmo",
		Amount: big.NewInt(1000000),
	}

	out, err := exchanger.Swap(in, []types.Pool{}, 980000)

	assert.Nil(t, err)

	_ = out
	_ = err
}
