package model

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
