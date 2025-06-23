package usecase

import (
	"context"
	"fmt"
	"log"
	"marketfuck/internal/application/port/out"
	"marketfuck/internal/domain/model"
	"marketfuck/internal/domain/service"
	"marketfuck/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// сохраняем данные с источника в кеш
func PriceAggregator(cache out.CacheClient, prices <-chan model.Price) {
	for price := range prices {
		price.TSR = time.Now().UnixMilli()

		key := fmt.Sprintf("%s:%s:%d", price.PairName, price.Exchange, price.TSR)

		// log.Printf("Добавляем цену в кеш: %s (TSR=%d, now=%d)", key, price.TSR, time.Now().UnixMilli())

		err := cache.SetPrice(context.Background(), key, price, 0)
		if err != nil {
			log.Printf("Ошибка установки цены в кеш: %v", err)
		}
	}
}

func CleanupOldPrices(cache out.CacheClient, thresholdMillis int64, delay int64) error {
	ctx := context.Background()
	keys, err := cache.Keys(ctx, "*")
	if err != nil {
		return fmt.Errorf("Ошибка получения ключей: %v", err)
	}
	for _, key := range keys {
		price, found, err := cache.GetPrice(ctx, key)
		if err != nil || !found {
			continue
		}

		if price.TSR <= thresholdMillis {
			err := cache.Delete(ctx, key)
			if err != nil {
				return fmt.Errorf("Ошибка при удалении ключа %s: %v", key, err)
			}
		}
	}
	return nil
}

func GetAllPrices(cache out.CacheClient, marketService *service.MarketService) ([]model.AggregatedPrice, int64) {
	ctx := context.Background()
	var recentKeys []string
	var cursor uint64 = 0
	maxIterations := 1000
	iteration := 0
	// возможно стоит передавать delay из main?
	from := time.Now().UnixMilli() - 10_001
	to := time.Now().UnixMilli()

	for {
		if iteration >= maxIterations {
			log.Printf("Достигнуто максимальное количество итераций: %d", maxIterations)
			break
		}

		keys, newCursor, err := cache.Scan(ctx, cursor, "*:*:*", 1000)
		if err != nil {
			log.Printf("Ошибка SCAN на итерации %d: %v", iteration, err)
			break
		}

		for _, key := range keys {
			parts := strings.Split(key, ":")
			if len(parts) != 3 {
				log.Printf("Ключ с неожиданной структурой: %s", key)
				continue
			}

			timestampMillis, err := strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				log.Printf("Ошибка парсинга timestamp из ключа %s: %v", key, err)
				continue
			}

			// Проверяем, попадает ли ключ в диапазон времени
			if timestampMillis >= from && timestampMillis <= to {
				// Получаем значение по ключу
				value, err := cache.Get(ctx, key)
				if err != nil {
					log.Printf("Ошибка получения значения для ключа %s: %v", key, err)
					continue
				}
				

				// Добавляем найденный ключ и его значение в результат
				recentKeys = append(recentKeys, fmt.Sprintf("%s: %s", key, value))
			}
		}

		cursor = newCursor
		iteration++

		// Завершаем сканирование, если курсор вернулся к нулю
		if cursor == 0 {
			log.Printf("Сканирование завершено (cursor = 0)")
			break
		}
	}

	count := len(recentKeys)
	log.Printf("Найдено ключей за последние 10 секунд: %d", count)
	aggr, err := utils.AggregatePricesByMinute(recentKeys)
	if err != nil {
		fmt.Println("некорректная агрегация", err)
	}
	return aggr, to
}
