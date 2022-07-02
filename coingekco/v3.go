package coingeckco

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/json"
)

const (
	rpcURL = "https://api.coingecko.com/api/v3"

	rpcTimeout = time.Second * 1
)

// Client is
type Client interface {
	Ping() error
	USDPrice(tickers []Ticker) (map[Ticker]float64, error)
}

type v3Client struct {
}

func (c *v3Client) Ping() error {
	req, err := http.NewRequest(http.MethodGet, rpcURL+"/ping", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	cli := http.Client{Timeout: rpcTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("no success poin response")
	}

	return nil
}

func (c *v3Client) USDPrice(tickers []Ticker) (map[Ticker]float64, error) {
	path := "/simple/price?vs_currencies=usd&ids="
	for _, ticker := range tickers {
		path += string(ticker) + ","
	}

	path = strings.TrimRight(path, ",")

	req, err := http.NewRequest(http.MethodGet, rpcURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	cli := http.Client{Timeout: rpcTimeout}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := map[Ticker]map[string]float64{}

	err = json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	retPrices := map[Ticker]float64{}

	for ticker, price := range respData {
		retPrices[ticker] = price["usd"]
	}

	return retPrices, nil
}

// NewClient create new CoinGecko client.
func NewClient() Client {
	return &v3Client{}
}
