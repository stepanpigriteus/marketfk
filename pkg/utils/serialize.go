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
	// Мапа для группировки
	grouped := make(map[string][]model.Price)

	// Разбираем строки в массив объектов Price
	for _, priceStr := range prices {
		// Разделяем строку на две части: до первого двоеточия и сам JSON
		parts := strings.SplitN(priceStr, ":", 4)
		if len(parts) != 4 {
			return nil, fmt.Errorf("неверный формат строки: %s", priceStr)
		}

		// Извлекаем символ (пару валют) и биржу
		pairName := parts[0]
		exchange := parts[1]
		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			log.Printf("Ошибка парсинга временной метки: %v", err)
			continue
		}

		// Разбираем JSON часть строки
		var priceData struct {
			Symbol   string  `json:"symbol"`
			Exchange string  `json:"Exchange"`
			Price    float64 `json:"Price"`
		}
		if err := json.Unmarshal([]byte(parts[3]), &priceData); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		// Используем timestamp для вычисления начала минуты
		minute := timestamp / 60000 * 60000 // округление до начала минуты
		groupKey := fmt.Sprintf("%s:%s:%d", pairName, exchange, minute)

		// Добавляем цену в группу
		grouped[groupKey] = append(grouped[groupKey], model.Price{
			PairName: pairName,
			Exchange: exchange,
			TSR:      timestamp,
			Price:    priceData.Price, // Получаем цену из JSON
		})
	}

	// Составляем итоговые агрегированные данные
	var results []model.AggregatedPrice
	for groupKey, prices := range grouped {
		// Разбираем ключ, чтобы извлечь PairName, Exchange и timestamp
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
