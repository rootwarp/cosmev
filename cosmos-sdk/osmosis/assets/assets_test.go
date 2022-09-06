package assets

import (
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAssetList(t *testing.T) {
	tests := []struct {
		Denom      string
		IsSuccess  bool
		ExpectUnit string
	}{
		{
			Denom:      "uosmo",
			IsSuccess:  true,
			ExpectUnit: "uosmo",
		},
		{
			Denom:      "uion",
			IsSuccess:  true,
			ExpectUnit: "uion",
		},
		{
			Denom:      "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
			IsSuccess:  true,
			ExpectUnit: "uatom",
		},
		{
			Denom:     "hello",
			IsSuccess: false,
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	fixture, err := os.ReadFile("./fixtures/assetlist.json")
	assert.Nil(t, err)

	httpmock.RegisterResponder(
		http.MethodGet,
		url,
		httpmock.NewStringResponder(http.StatusOK, string(fixture)))

	cli := NewClient()

	err = cli.Fetch()
	assert.Nil(t, err)

	for _, test := range tests {
		denom, err := cli.Denom(test.Denom)

		if test.IsSuccess {
			assert.Nil(t, err)
			assert.Equal(t, test.ExpectUnit, denom.Unit())
		} else {
			assert.NotNil(t, err)
		}
	}
}
