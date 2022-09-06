package types

import (
	"math/big"
)

// Pool is swap pool information of Osmosis.
type Pool struct {
	ID         string
	Address    string
	SwapFee    float64
	PoolAssets []PoolAsset
}

// PoolAsset defines detail asset informations from pool.
type PoolAsset struct {
	Denom  string
	Amount *big.Int
	Weight *big.Int
}
