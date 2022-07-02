package coingeckco

// Ticker describes unique ID of token from Coingecko.
type Ticker string

// Tickers from Coingecko
const (
	Atom      Ticker = "cosmos"
	Osmosis          = "osmosis"
	Ion              = "ion"
	Evmos            = "evmos"
	Chihuahua        = "chihuahua-token"
	WETH             = "weth"
	MediBloc         = "medibloc"

	USDC Ticker = "usd-coin"
	DAI         = "dai"
)
