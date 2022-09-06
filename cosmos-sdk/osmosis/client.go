package osmosis

import (
	osmoRpc "github.com/rootwarp/cosmev/cosmos-sdk/osmosis/rpc"
	"github.com/rootwarp/cosmev/types"
)

var (
	//url = "https://osmosis-mainnet-rpc.allthatnode.com:1317"
	url = "https://osmosis.stakesystems.io"
)

// NewClient creates DexReader for Osmosis.
func NewClient() types.DexReader {
	return osmoRpc.NewPoolClient(url)
}
