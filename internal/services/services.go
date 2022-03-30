package services

import (
	"context"
	"time"
)

type RankService struct {
}

func NewRankService() *RankService {
	return &RankService{}
}

type PriceService struct {
}

func NewPriceService() *PriceService {
	return &PriceService{}
}

type ApiService struct {
}

func NewApiService() *ApiService {
	return &ApiService{}
}

func (a *ApiService) Hello(ctx context.Context) (res context.Context, err error) {
	time.Sleep(10 * time.Second)

	return ctx, nil
}
