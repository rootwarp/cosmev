package coingeckco

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	cli := NewClient()
	err := cli.Ping()

	assert.Nil(t, err)
}

func TestPrice(t *testing.T) {
	cli := NewClient()

	reqTickers := []Ticker{Osmosis, Ion, USDC, DAI, Chihuahua, WETH, MediBloc}

	prices, err := cli.USDPrice(reqTickers)
	assert.Nil(t, err)

	for _, ticker := range reqTickers {
		_, ok := prices[ticker]
		assert.True(t, ok)
	}
}
