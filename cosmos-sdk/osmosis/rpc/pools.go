package osmosis

import (
	"io"
	"net/http"

	"github.com/tendermint/tendermint/libs/json"
)

type listPoolResp struct {
	Pools      []Pool `json:"pools"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}

// Pool is ...
type Pool struct {
	Type       string `json:"@type"`
	Address    string `json:"address"`
	ID         string `json:"id"`
	PoolParams struct {
		SwapFee                  string `json:"swapFee"`
		ExitFee                  string `json:"exitFee"`
		SmoothWeightChangeParams string `json:"smoothWeightChangeParams"`
	} `json:"poolParams"`
	FuturePoolGovernor string `json:"future_pool_governor"`
	TotalShares        struct {
		Demon  string `json:"demon"`
		Amount string `json:"amount"`
	} `json:"totalShares"`
	PoolAssets []struct {
		Token struct {
			Demon  string `json:"demon"`
			Amount string `json:"amount"`
		} `json:"token"`
		Weight string `json:"weight"`
	} `json:"poolAssets"`
	TotalWeight string `json:"totalWeight"`
}

// PoolReader is ...
type PoolReader interface {
	List() (map[string]Pool, error)
}

type poolReader struct {
	rpcURL string
}

func (r *poolReader) List() (map[string]Pool, error) {
	pools := map[string]Pool{}

	lastPaginationKey := ""
	for {
		url := r.rpcURL + "/osmosis/gamm/v1beta1/pools?pagination.limit=1000"
		if len(pools) > 0 {
			url += "&pagination.key=" + lastPaginationKey
		}

		data, err := httpQuery(url)
		if err != nil {
			return nil, err
		}

		respData := new(listPoolResp)
		err = json.Unmarshal(data, respData)
		if err != nil {
			return nil, err
		}

		lastPaginationKey = respData.Pagination.NextKey
		for _, p := range respData.Pools {
			pools[p.ID] = p
		}

		if respData.Pagination.NextKey == "" {
			break
		}
	}

	return pools, nil
}

func httpQuery(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	rawBody, err := io.ReadAll(resp.Body)
	return rawBody, err
}

// NewPoolClient creates instance of Pool.
func NewPoolClient(rpcURL string) PoolReader {
	return &poolReader{
		rpcURL: rpcURL,
	}
}
