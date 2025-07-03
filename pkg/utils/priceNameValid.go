package utils

import (
	"fmt"
	"strings"
)

func PairNameValidFormatter(pairName string) (string, error) {
	validPairs := map[string]bool{
		"BTCUSDT":  true,
		"DOGEUSDT": true,
		"TONUSDT":  true,
		"SOLUSDT":  true,
		"ETHUSDT":  true,
	}
	pairName = strings.ToUpper(pairName)
	if len(pairName) < 7 || !validPairs[pairName] {
		return "", fmt.Errorf("incorrect PairName")
	}
	return pairName, nil
}
