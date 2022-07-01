package osmosis

import (
	"encoding/json"
	"io"
	"net/http"

	cosmossdk "github.com/rootwarp/cosmev/cosmos-sdk"
)

type accountResponse struct {
	Account cosmossdk.Account `json:"account"`
}

type auth struct {
	rpcURL string
}

func (a *auth) GetAccount(address string) (*cosmossdk.Account, error) {
	req, err := http.NewRequest(http.MethodGet, a.rpcURL+"/cosmos/auth/v1beta1/accounts/"+address, nil)
	if err != nil {
		return nil, err
	}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	accResp := accountResponse{}
	err = json.Unmarshal(data, &accResp)
	if err != nil {
		return nil, err
	}

	return &accResp.Account, nil
}

// NewAuthClient creates new client of Auth module.
func NewAuthClient(rpcURL string) cosmossdk.Auth {
	return &auth{
		rpcURL: rpcURL,
	}
}
