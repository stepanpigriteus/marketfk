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

		// minute := timestamp / 60000 * 60000
		groupKey := fmt.Sprintf("%s:%s:%d", pairName, exchange, timestamp)

		grouped[groupKey] = append(grouped[groupKey], model.Price{
			PairName: pairName,
			Exchange: exchange,
			TSR:      timestamp,
			Price:    priceData.Price,
		})
		// log.Printf("Добавлена цена: Pair=%s, Exchange=%s, Timestamp=%d, Price=%.2f", pairName, exchange, timestamp, priceData.Price)
	}

	var results []model.AggregatedPrice
	for groupKey, prices := range grouped {
		parts := strings.Split(groupKey, ":")
		pair, exchange := parts[0], parts[1]
		minuteTs, _ := strconv.ParseInt(parts[2], 10, 64)

		var sum, min, max float64
		count := len(prices)

		// log.Printf("Обработка группы: %s, Количество цен: %d", groupKey, count)

		if count > 0 {
			min = prices[0].Price
			max = prices[0].Price
			sum = prices[0].Price

			// log.Printf("Начальные значения: min=%.2f, max=%.2f, sum=%.2f", min, max, sum)

			for i := 1; i < count; i++ {
				price := prices[i].Price
				sum += price
				if price < min {
					min = price
				}
				if price > max {
					max = price
				}
				// log.Printf("Цена %d: %.2f, min=%.2f, max=%.2f, sum=%.2f", i, price, min, max, sum)
			}

			avg := sum / float64(count)

			results = append(results, model.AggregatedPrice{
				PairName:     pair,
				Exchange:     exchange,
				Timestamp:    time.UnixMilli(minuteTs),
				AveragePrice: avg,
				MinPrice:     min,
				MaxPrice:     max,
			})
			// log.Printf("Результат для группы %s: avg=%.2f, min=%.2f, max=%.2f", groupKey, avg, min, max)
		} else {
			// log.Printf("Пустая группа: %s", groupKey)
		}
	}

	return results, nil
}
