package osmosis

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/rootwarp/cosmev/types"
	"github.com/stretchr/testify/assert"
)

var (
	testMnemonic = os.Getenv("COSMEV_TEST_MNEMONIC")
	testAddress  = os.Getenv("COSMEV_TEST_ADDRESS")
)

func TestExchange(t *testing.T) {
	exchanger := NewExchanger(os.Getenv("COSMEV_TEST_RPC"))

	hdParam := hd.CreateHDPath(118, 0, 0)
	err := exchanger.SetMnemonic(testMnemonic, hdParam.String())

	assert.Nil(t, err)
	assert.Equal(t, testAddress, exchanger.Address())

	// TODO: Find path and add.
	// TODO: Check pool route is valid.
	// Pool -> SwapMountInRoute

	in := types.Asset{
		Denom:  "uosmo",
		Amount: big.NewInt(1000000),
	}

	out, err := exchanger.Swap(in, []types.Pool{}, 980000)

	assert.Nil(t, err)

	_ = out
	_ = err
}

func TestRoutePath(t *testing.T) {
	/*
		Example
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
	*/

	exchanger := osmosisExchanger{}

	path := []types.Pool{
		{
			ID:      "1",
			Address: "osmo1mw0ac6rwlp5r8wapwk3zs6g29h8fcscxqakdzw9emkne6c8wjp9q0t3v8t",
			SwapFee: 0.002,
			PoolAssets: []types.PoolAsset{
				{Denom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"},
				{Denom: "uosmo"},
			},
		},

		{
			ID:      "498",
			Address: "osmo1tusadtwjnzzyakm94t5gjqr4dlkdcp63hctlql6xvslvkf7kkdws5lfyxc",
			SwapFee: 0.003,
			PoolAssets: []types.PoolAsset{
				{Denom: "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"},
				{Denom: "ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED"},
			},
		},

		{
			ID:      "497",
			Address: "osmo1h7yfu7x4qsv2urnkl4kzydgxegdfyjdry5ee4xzj98jwz0uh07rqdkmprr",
			SwapFee: 0.003,
			PoolAssets: []types.PoolAsset{
				{Denom: "ibc/46B44899322F3CD854D2D46DEEF881958467CDD4B3B10086DA49296BBED94BED"},
				{Denom: "uosmo"},
			},
		},
	}

	out, err := exchanger.convert(path, "uosmo")

	fmt.Println(out, err)
}
