package models

type Coin struct {
	NameCoin string `json:"name_coin" binding:"required"`
}

type CoinPrice struct {
	USD float64 `json:"usd"`
}
