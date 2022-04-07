package provider

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/StepanchukYI/top-coin/internal/config"
	"github.com/StepanchukYI/top-coin/internal/models"
)

type RankProvider struct {
	Url      string
	ApiKey   string
	Currency string
}

func NewRankProvider(config *config.Config) *RankProvider {
	return &RankProvider{
		ApiKey:   config.CryptocompareApiKey,
		Url:      config.CryptocompareApiUrl,
		Currency: config.Currency,
	}
}

func (p *RankProvider) GetRank(limit int, page int) ([]models.Crypto, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.Url, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	q.Add("tsym", p.Currency)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("authorization", "Apikey "+p.ApiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error sending request to server")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error receiving request to server from Rank Data Provider")
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	jsonRes := new(coinsRoot)
	err = json.Unmarshal(respBody, &jsonRes)
	if err != nil {
		return nil, err
	}

	var Cryptos []models.Crypto
	for key, value := range jsonRes.Data {
		coin, ok := value["CoinInfo"]
		if !ok {
			return nil, errors.New("CoinInfo not found into API responce ")
		}

		symbol := coin.Name
		Cryptos = append(Cryptos, models.Crypto{
			Rank:   (page * limit) + key + 1,
			Symbol: symbol,
		})
	}

	return Cryptos, nil
}

type coinsRoot struct {
	Data []map[string]Coin `json:"Data"`
}

type Coin struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}
