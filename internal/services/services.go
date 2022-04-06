package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/StepanchukYI/top-coin/internal/models"
	"github.com/StepanchukYI/top-coin/internal/provider"
)

type RankService struct {
	RankProvider *provider.RankProvider
}

func NewRankService(rank *provider.RankProvider) *RankService {
	return &RankService{
		RankProvider: rank,
	}
}

func (s *RankService) Rank(ctx context.Context) ([]models.Crypto, error) {
	defaultLimit := 20
	maxLimit := 100
	limit := defaultLimit

	limitVal := ctx.Value("limit")

	if limitVal != nil {
		limit = limitVal.(int)
	}
	if limit > maxLimit {
		return nil, errors.New("Limit must be lowest than " + string(maxLimit))
	}
	
	data, err := s.RankProvider.GetRank(limit)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type PriceService struct {
	PriceProvider *provider.PriceProvider
}

func NewPriceService(price *provider.PriceProvider) *PriceService {
	return &PriceService{
		PriceProvider: price,
	}
}

type ApiService struct {
	RankProvider  *provider.RankProvider
	PriceProvider *provider.PriceProvider
}

func NewApiService(rank *provider.RankProvider, price *provider.PriceProvider) *ApiService {
	return &ApiService{
		RankProvider:  rank,
		PriceProvider: price,
	}
}

func (a *ApiService) Currency(ctx context.Context) ([]models.Crypto, error) {
	defaultLimit := 20
	maxLimit := 100
	limit := defaultLimit

	limitVal := ctx.Value("limit")

	if limitVal != nil {
		limit = limitVal.(int)
	}
	if limit > maxLimit {
		return nil, errors.New("Limit must be lowest than " + string(maxLimit))
	}

	cryptos, err := a.RankProvider.GetRank(limit)
	if err != nil {
		return nil, err
	}

	currencyKeys := []string{}

	for _, crypto := range cryptos {
		currencyKeys = append(currencyKeys, crypto.Symbol)
	}

	prices, err := a.PriceProvider.GetPrice(strings.Join(currencyKeys, ","))

	if err != nil {
		return nil, err
	}

	for _, crypto := range cryptos {
		cryptoPrice, ok := prices[crypto.Symbol]
		if !ok {

		}
		crypto.Price = cryptoPrice
	}

	return cryptos, nil
}

func (a *ApiService) Hello(ctx context.Context) (interface{}, error) {

	time.Sleep(5 * time.Second)

	random := rand.Intn(5)

	fmt.Println(random)

	if random == 1 {
		return ctx, nil
	} else {
		return ctx, errors.New("Error")
	}

}
