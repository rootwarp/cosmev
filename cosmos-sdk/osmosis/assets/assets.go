package assets

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	// Check here: https://github.com/osmosis-labs/assetlists/blob/main/osmosis-1/osmosis-1.assetlist.json
	url = "https://raw.githubusercontent.com/osmosis-labs/assetlists/main/osmosis-1/osmosis-1.assetlist.json"
)

// AssetList defines Osmosis DEX list.
type AssetList struct {
	ChainID string  `json:"chain_id"`
	Assets  []Asset `json:"assets"`
}

// Asset defines detail information of listed asset.
type Asset struct {
	Description string      `json:"description"`
	DenomUnits  []DenomUnit `json:"denom_units"`
	Base        string      `json:"base"`
	Name        string      `json:"name"`
	Display     string      `json:"display"`
	Symbol      string      `json:"symbol"`
	IBC         IBC         `json:"ibc"`
	CoingeckoID string      `json:"coingecko_id"`
}

// DenomUnit defines unit information.
type DenomUnit struct {
	Denom    string   `json:"denom"`
	Exponent int      `json:"exponent"`
	Aliases  []string `json:"aliases"`
}

// Unit returns ticker name of denomunit.
func (d DenomUnit) Unit() string {
	if len(d.Aliases) != 0 {
		return d.Aliases[0]
	}

	return d.Denom
}

// IBC defines ibc channel info.
type IBC struct {
	SourceChannel string `json:"source_channel"`
	DstChannel    string `json:"dst_channel"`
	SourceDenom   string `json:"source_denom"`
}

// AssetListReader defines getter interfaces for Osmosis DEX list.
type AssetListReader interface {
	Fetch() error
	Denom(denom string) (*DenomUnit, error)
}

type assetListClient struct {
	assets map[string]Asset
}

func (c *assetListClient) Fetch() error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	l := new(AssetList)
	err = json.Unmarshal(rawData, l)
	if err != nil {
		return err
	}

	for _, asset := range l.Assets {
		c.assets[asset.Base] = asset
	}

	return nil
}

func (c *assetListClient) Denom(denom string) (*DenomUnit, error) {
	asset, ok := c.assets[denom]
	if !ok {
		return nil, errors.New("cannot find denom")
	}

	var findUnit *DenomUnit
	for _, unit := range asset.DenomUnits {
		if unit.Denom == denom {
			findUnit = &unit
			break
		}
	}

	return findUnit, nil
}

// NewClient creates new reader client.
func NewClient() AssetListReader {
	return &assetListClient{
		assets: map[string]Asset{},
	}
}
