package osmosis

import (
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/rootwarp/cosmev/types"
	"github.com/tendermint/tendermint/libs/json"
)

type listPoolResp struct {
	Pools      []pool `json:"pools"`
	Pagination struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	} `json:"pagination"`
}

type pool struct {
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
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"totalShares"`
	PoolAssets []struct {
		Token struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"token"`
		Weight string `json:"weight"`
	} `json:"poolAssets"`
	TotalWeight string `json:"totalWeight"`
}

func (p pool) Convert() (*types.Pool, error) {
	newPool := types.Pool{
		ID:         p.ID,
		Address:    p.Address,
		PoolAssets: make([]types.PoolAsset, len(p.PoolAssets)),
	}

	swapFee, err := strconv.ParseFloat(p.PoolParams.SwapFee, 64)
	if err != nil {
		return nil, err
	}

	newPool.SwapFee = swapFee

	for i, asset := range p.PoolAssets {
		newPoolAsset := types.PoolAsset{}
		newPoolAsset.Denom = asset.Token.Denom

		amount, ok := new(big.Int).SetString(asset.Token.Amount, 10)
		if !ok {
			return nil, errors.New("can't parse amount value")
		}

		newPoolAsset.Amount = amount

		weight, ok := new(big.Int).SetString(asset.Weight, 10)
		if !ok {
			return nil, errors.New("can't parse weight value")
		}

		newPoolAsset.Weight = weight

		newPool.PoolAssets[i] = newPoolAsset
	}

	return &newPool, nil
}

// PoolReader defines interfaces for Pool.
type PoolReader interface {
	List() (map[string]*types.Pool, error)
}

type poolReader struct {
	rpcURL string
}

func (r *poolReader) List() (map[string]*types.Pool, error) {
	pools := map[string]*types.Pool{}

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
			parsePool, err := p.Convert()
			if err != nil {
				log.Println(err)
				continue
			}

			pools[p.ID] = parsePool
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
