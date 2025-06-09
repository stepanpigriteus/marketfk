package model

import "time"

type Price struct {
	PairName  string `json:"symbol"`
	Exchange  string
	Price     float64
	Timestamp int64 `json:"timestamp"`
	TSR       int64 `json:"timestamp"`
}

type Pair struct {
	Name string
}

type AggregatedPrice struct {
	PairName     string
	Exchange     string
	Timestamp    time.Time
	AveragePrice float64
	MinPrice     float64
	MaxPrice     float64
}
