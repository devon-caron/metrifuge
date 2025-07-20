package types

// Coin Balance Params
type CoinBalanceParams struct {
	Username string
}

// Coin Balance Response
type CoinBalanceResponse struct {
	// HTTP Status Code
	Code int

	// Account Balance
	Balance int64
}
