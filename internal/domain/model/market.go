package model

import "time"


type Price struct {
	PairName  string  
	Exchange  string   
	Price     float64   
	Timestamp time.Time
}


type Pair struct {
	Name string 
}
