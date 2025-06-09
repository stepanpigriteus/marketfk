package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"marketfuck/internal/domain/model"
	"strconv"
	"strings"
	"time"
)

func AggregatePricesByMinute(prices []string) ([]model.AggregatedPrice, error) {
	grouped := make(map[string][]model.Price)

	for _, priceStr := range prices {

		parts := strings.SplitN(priceStr, ":", 4)
		if len(parts) != 4 {
			return nil, fmt.Errorf("неверный формат строки: %s", priceStr)
		}

		pairName := parts[0]
		exchange := parts[1]
		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			log.Printf("Ошибка парсинга временной метки: %v", err)
			continue
		}

		var priceData struct {
			Symbol   string  `json:"symbol"`
			Exchange string  `json:"Exchange"`
			Price    float64 `json:"Price"`
		}
		if err := json.Unmarshal([]byte(parts[3]), &priceData); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		minute := timestamp / 60000 * 60000
		groupKey := fmt.Sprintf("%s:%s:%d", pairName, exchange, minute)

		grouped[groupKey] = append(grouped[groupKey], model.Price{
			PairName: pairName,
			Exchange: exchange,
			TSR:      timestamp,
			Price:    priceData.Price,
		})
	}

	var results []model.AggregatedPrice
	for groupKey, prices := range grouped {

		parts := strings.Split(groupKey, ":")
		pair, exchange := parts[0], parts[1]
		minuteTs, _ := strconv.ParseInt(parts[2], 10, 64)

		var sum, min, max float64
		for i, p := range prices {
			if i == 0 {
				min, max = p.Price, p.Price
			} else {
				if p.Price < min {
					min = p.Price
				}
				if p.Price > max {
					max = p.Price
				}
			}
			sum += p.Price
		}

		avg := sum / float64(len(prices))

		results = append(results, model.AggregatedPrice{
			PairName:     pair,
			Exchange:     exchange,
			Timestamp:    time.UnixMilli(minuteTs),
			AveragePrice: avg,
			MinPrice:     min,
			MaxPrice:     max,
		})
	}

	return results, nil
}
