package models

type PriceRequest struct {
	Coin      string `json:"coin" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

type DbResponse struct {
	Coin      string  `db:"coin"`
	Price     float64 `db:"price"`
	Precision int     `db:"precision"`
	Currency  string  `db:"currency"`
	Timestamp int64   `db:"timestamp"`
}

type PriceResponse struct {
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	Timestamp int64   `json:"timestamp"`
}
