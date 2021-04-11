package domain

type Asset struct {
	Ticker string  `json:"ticker"`
	Amount float64 `json:"amount"`
}
