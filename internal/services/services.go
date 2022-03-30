package services

import (
	"context"
	"fmt"
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

func (a *ApiService) Hello(ctx context.Context) (interface{}, error) {

	fmt.Println("1")
	time.Sleep(5 * time.Second)

	return ctx, nil
}
