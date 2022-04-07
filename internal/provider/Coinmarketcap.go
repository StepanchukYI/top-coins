package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/StepanchukYI/top-coin/internal/config"
)

type PriceProvider struct {
	Url      string
	ApiKey   string
	Currency string
}

func NewPriceProvider(config *config.Config) *PriceProvider {
	return &PriceProvider{
		ApiKey:   config.CoinmarketcapApiKey,
		Url:      config.CoinmarketcapApiUrl,
		Currency: config.Currency,
	}
}

func (p *PriceProvider) GetPrice(crypto string) (map[string]float64, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.Url, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(crypto)
	q := url.Values{}
	q.Add("symbol", crypto)
	q.Add("aux", "is_active")
	q.Add("convert", p.Currency)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", p.ApiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error sending request to server")
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%+v\n", resp)
		fmt.Printf("%+v\n", resp.Body)
		return nil, errors.New("Error receiving request to server from Price Data Provider")
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	jsonRes := new(marketRoot)
	err = json.Unmarshal(respBody, &jsonRes)
	if err != nil {
		return nil, err
	}

	response := make(map[string]float64)

	for k, v := range jsonRes.Data {

		price, ok := v.Quote[p.Currency]
		if !ok {
			return nil, errors.New("Price not found into API responce ")
		}

		response[k] = float64(price.Price)
	}

	return response, nil
}

type marketRoot struct {
	Data map[string]Quote `json:"Data"`
}

type Quote struct {
	Quote map[string]CoinPrice `json:"quote"`
}

type CoinPrice struct {
	Price float64 `json:"price"`
}
