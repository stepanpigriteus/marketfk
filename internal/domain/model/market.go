package model

type Price struct {
	PairName  string
	Exchange  string
	Price     float64
	Timestamp int64 `json:"timestamp"`
}

type Pair struct {
	Name string
}
