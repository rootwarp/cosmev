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

// Asset is token balance of account.
type Asset struct {
	Denom  string
	Amount *big.Int
}

// DexReader defines interfaces for DEX.
type DexReader interface {
	ListPool() (map[string]*Pool, error)
}

// DexExchanger defines interfaces for Swap.
type DexExchanger interface {
	SetMnemonic(mnemonic, hdPath string) error
	Address() string
	Swap(tokenIn Asset, routes []Pool, minTokenOutAmount int64) (*Asset, error)
}

// Account
// Fee
// Gaslimit
