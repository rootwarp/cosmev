package osmosis

import (
	"fmt"
	"os"
	"testing"

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
	r := NewPoolClient(rpcURL)

	pools, err := r.List()

	fmt.Println(len(pools), err)
	fmt.Printf("%+v\n", pools["1"])

	assert.Nil(t, nil)
}
