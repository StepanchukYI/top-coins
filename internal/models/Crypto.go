package models

type Crypto struct {
	Rank   int     `json:"rank"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}
