package osmosis

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var (
	rpcURL = os.Getenv("COSMEV_TEST_RPC")
)

func TestNumPools(t *testing.T) {
	data, err := httpQuery(rpcURL + "/osmosis/gamm/v1beta1/num_pools")
	fmt.Println(string(data), err)

	assert.Nil(t, nil)
}

func TestListPool(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixture, err := os.ReadFile("./fixtures/pools.txt")
	assert.Nil(t, err)

	rpcURL := "https://test.dev"
	httpmock.RegisterResponder(
		http.MethodGet,
		rpcURL+"/osmosis/gamm/v1beta1/pools?pagination.limit=1000",
		httpmock.NewStringResponder(http.StatusOK, string(fixture)))

	r := NewPoolClient(rpcURL)

	pools, err := r.List()

	assert.Nil(t, err)

	pool1, ok := pools["1"]

	assert.True(t, ok)
	assert.Equal(t, "1", pool1.ID)
	assert.Equal(t, "osmo1mw0ac6rwlp5r8wapwk3zs6g29h8fcscxqakdzw9emkne6c8wjp9q0t3v8t", pool1.Address)
	assert.Equal(t, 0.002, pool1.SwapFee)

	assert.Equal(t, "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", pool1.PoolAssets[0].Denom)
	assert.Equal(t, "3370086542043", pool1.PoolAssets[0].Amount.String())
	assert.Equal(t, "536870912000000", pool1.PoolAssets[0].Weight.String())

	assert.Equal(t, "uosmo", pool1.PoolAssets[1].Denom)
	assert.Equal(t, "35008926848057", pool1.PoolAssets[1].Amount.String())
	assert.Equal(t, "536870912000000", pool1.PoolAssets[1].Weight.String())
}
