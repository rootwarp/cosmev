package tests

import (
	"testing"

	"github.com/rootwarp/cosmev/cosmos-sdk/osmosis"
	"github.com/stretchr/testify/assert"
)

func TestOsmosisPool(t *testing.T) {
	r := osmosis.NewClient()

	poolMap, err := r.ListPool()

	assert.Nil(t, err)

	pool1, ok := poolMap["1"]

	assert.True(t, ok)
	assert.Equal(t, "uatom", pool1.PoolAssets[0].Denom)
	assert.Equal(t, "uosmo", pool1.PoolAssets[1].Denom)
}
