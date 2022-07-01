package cosmossdk

// Account is user's account information.
type Account struct {
	Type    string `json:"@type"`
	Address string `json:"address"`
	Pubkey  struct {
		Type string `json:"@type"`
		Key  string `json:"key"`
	} `json:"pub_key"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

// Auth provides interfaces to handle account
type Auth interface {
	GetAccount(address string) (*Account, error)
}
