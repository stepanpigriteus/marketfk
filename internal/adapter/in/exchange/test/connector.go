package test

import (
	"marketfuck/internal/domain/model"
	"sync"
	"time"
)

type TestGenerator struct {
	Exchange     model.Exchange
	Pairs        []string
	BasePrice    map[string]float64
	PriceChanges map[string]float64
	Ticker       *time.Ticker
	Mu           sync.RWMutex
}

func NewTestGenerator(exchangeName string) *TestGenerator {
	pairs := []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	basePrice := map[string]float64{
		"BTCUSDT":  45000.0,
		"DOGEUSDT": 0.08,
		"TONUSDT":  2.5,
		"SOLUSDT":  100.0,
		"ETHUSDT":  3000.0,
	}

	priceChanges := map[string]float64{
		"BTCUSDT":  0.01,  // 1%
		"DOGEUSDT": 0.05,  // 5%
		"TONUSDT":  0.03,  // 3%
		"SOLUSDT":  0.02,  // 2%
		"ETHUSDT":  0.015, // 1.5%
	}

	return &TestGenerator{
		Exchange:     model.Exchange{Name: exchangeName},
		Pairs:        pairs,
		BasePrice:    basePrice,
		PriceChanges: priceChanges,
		Ticker:       time.NewTicker(100 * time.Millisecond), // Генерация каждые 100мс
	}
}
