package services

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/StepanchukYI/top-coin/internal/models"
	"github.com/StepanchukYI/top-coin/internal/provider"
	"github.com/StepanchukYI/top-coin/internal/server"
)

type RankService struct {
	RankProvider *provider.RankProvider
}

func NewRankService(rank *provider.RankProvider) *RankService {
	return &RankService{
		RankProvider: rank,
	}
}

func (s *RankService) Rank(ctx context.Context) ([]models.Crypto, server.ErrorResponse) {
	defaultLimit := 20
	maxLimit := 100
	limit := defaultLimit
	err := errors.New("")

	limitVal := ctx.Value("limit")
	if limitVal != nil {
		limit, err = strconv.Atoi(limitVal.(string))
		if err != nil {
			return nil, server.ErrorResponse{
				Code:   http.StatusInternalServerError,
				Errors: []string{err.Error()},
			}
		}
	}

	if limit%maxLimit != 0 {
		return nil, server.ErrorResponse{
			Code:   http.StatusInternalServerError,
			Errors: []string{"Limit must equal 100"},
		}
	}

	limitRequesrs := []int{}

	currentLimit := limit
	for {
		currentLimit -= maxLimit
		limitRequesrs = append(limitRequesrs, maxLimit)
		if currentLimit == maxLimit {
			limitRequesrs = append(limitRequesrs, maxLimit)
			break
		}
	}

	wg := &sync.WaitGroup{}
	responseModels := []models.Crypto{}

	for page, limit := range limitRequesrs {
		wg.Add(1)
		data, err := s.RankProvider.GetRank(limit, page, wg)
		if err != nil {
			return nil, server.ErrorResponse{
				Code:   http.StatusInternalServerError,
				Errors: []string{err.Error()},
			}
		}
		responseModels = append(responseModels, data...)
	}

	wg.Wait()

	return responseModels, server.ErrorResponse{}
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

func (a *ApiService) Currency(ctx context.Context) ([]models.Crypto, server.ErrorResponse) {
	defaultLimit := 20
	maxLimit := 100
	limit := defaultLimit
	err := errors.New("")

	limitVal := ctx.Value("limit")
	if limitVal != nil {
		limit, err = strconv.Atoi(limitVal.(string))
		if err != nil {
			return nil, server.ErrorResponse{
				Code:   http.StatusInternalServerError,
				Errors: []string{err.Error()},
			}
		}
	}

	if limit%maxLimit != 0 {
		return nil, server.ErrorResponse{
			Code:   http.StatusInternalServerError,
			Errors: []string{"Limit must equal 100"},
		}
	}

	limitRequesrs := []int{}

	currentLimit := limit
	for {
		currentLimit -= maxLimit
		limitRequesrs = append(limitRequesrs, maxLimit)
		if currentLimit == maxLimit {
			limitRequesrs = append(limitRequesrs, maxLimit)
			break
		}
	}

	wg := &sync.WaitGroup{}
	responseModels := []models.Crypto{}

	for page, limit := range limitRequesrs {
		wg.Add(1)
		data, err := a.RankProvider.GetRank(limit, page, wg)
		if err != nil {
			return nil, server.ErrorResponse{
				Code:   http.StatusInternalServerError,
				Errors: []string{err.Error()},
			}
		}
		responseModels = append(responseModels, data...)
	}
	wg.Wait()

	currencyKeys := []string{}

	for _, crypto := range responseModels {
		currencyKeys = append(currencyKeys, crypto.Symbol)
	}

	prices, err := a.PriceProvider.GetPrice(strings.Join(currencyKeys, ","))

	if err != nil {
		return nil, server.ErrorResponse{
			Code:   http.StatusInternalServerError,
			Errors: []string{err.Error()},
		}
	}

	for key, crypto := range responseModels {
		cryptoPrice, ok := prices[crypto.Symbol]
		if !ok {
			return nil, server.ErrorResponse{
				Code:   http.StatusBadRequest,
				Errors: []string{"API ERROR"},
			}
		}
		err = crypto.SetPrice(cryptoPrice)
		if err != nil {
			return nil, server.ErrorResponse{
				Code:   http.StatusInternalServerError,
				Errors: []string{err.Error()},
			}
		}
		responseModels[key] = crypto
	}

	return responseModels, server.ErrorResponse{}
}
