package services

import (
	"context"
	"errors"
	"fmt"
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
	maxLimit := 100
	limit := maxLimit
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
		if currentLimit == 0 {
			break
		}
	}

	responseModels := []models.Crypto{}
	errorsChan := make(chan server.ErrorResponse)

	for page, limit := range limitRequesrs {
		data, err := s.RankProvider.GetRank(limit, page)
		if err != nil {
			errorsChan <- server.ErrorResponse{
				Code:   http.StatusInternalServerError,
				Errors: []string{err.Error()},
			}
		}
		responseModels = append(responseModels, data...)
	}

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
	maxLimit := 100
	limit := maxLimit
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

	var wg sync.WaitGroup

	responseModels := []models.Crypto{}
	responseChan := make(chan []models.Crypto)
	errorsChan := make(chan server.ErrorResponse)

	for page, limit := range limitRequesrs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := a.RankProvider.GetRank(limit, page)
			if err != nil {
				errorsChan <- server.ErrorResponse{
					Code:   http.StatusInternalServerError,
					Errors: []string{err.Error()},
				}
			}
			fmt.Println("Request Done")
			responseChan <- data
		}()
	}
	fmt.Println("Wait")
	wg.Wait()
	fmt.Println("Wait ended")

	for data := range responseChan {
		responseModels = append(responseModels, data...)
	}

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
