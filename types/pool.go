package types

import (
	"fmt"
	"math/big"
)

// Pool is swap pool information of Osmosis.
type Pool struct {
	ID         string
	Address    string
	SwapFee    float64
	PoolAssets []PoolAsset
}

// FindPoolAssetByDenom searchs request denom and reset PoolAsset.
func (p Pool) FindPoolAssetByDenom(denom string) (int, *PoolAsset, error) {
	for i, asset := range p.PoolAssets {
		if asset.Denom == denom {
			return i, &asset, nil
		}
	}

	return -1, nil, fmt.Errorf("asset %s is not exist in this pool", denom)
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
