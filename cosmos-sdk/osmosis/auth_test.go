package osmosis

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	rpc     = os.Getenv("COSMEV_TEST_RPC")
	address = os.Getenv("COSMEV_TEST_ADDRESS")
)

func init() {
	cfg := types.GetConfig()
	cfg.SetBech32PrefixForAccount("osmo", "osmopub")
	cfg.SetBech32PrefixForValidator("osmovaloper", "osmovaloperpub")
	cfg.SetBech32PrefixForConsensusNode("osmovalcons", "osmovalconspub")
	cfg.Seal()
}

func TestGetAccount(t *testing.T) {
	cli := NewAuthClient(rpc)
	acc, err := cli.GetAccount(address)

	assert.Nil(t, err)
	assert.Equal(t, address, acc.Address)
}
