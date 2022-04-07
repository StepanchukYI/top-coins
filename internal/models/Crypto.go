package models

type Crypto struct {
	Rank   int     `json:"Rank"`
	Symbol string  `json:"Symbol"`
	Price  float64 `json:"Price USD"`
}

func (c *Crypto) SetPrice(price float64) error {
	c.Price = price
	return nil
}
