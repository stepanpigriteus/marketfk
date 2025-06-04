package usecase

import (
	"context"
	"fmt"
	"log"
	"marketfuck/internal/application/port/out"
	"marketfuck/internal/domain/model"
	"time"
)

func PriceAggregator(cache out.CacheClient, prices <-chan model.Price) {
	count := 0
	for price := range prices {

		price.TSR = time.Now().UnixMilli()

		key := fmt.Sprintf("%s:%s:%d", price.PairName, price.Exchange, price.TSR)

		err := cache.SetPrice(context.Background(), key, price, 0)
		if err != nil {
			log.Printf("Ошибка установки цены в кеш: %v", err)
		}
		count++
		fmt.Println(count)
	}
}

func CleanupOldPrices(cache out.CacheClient, thresholdMillis int64) {
	ctx := context.Background()
	keys, err := cache.Keys(ctx, "*")
	if err != nil {
		log.Printf("Ошибка получения ключей: %v", err)
		return
	}

	now := time.Now().UnixMilli()
	for _, key := range keys {
		price, found, err := cache.GetPrice(ctx, key)
		if err != nil || !found {
			continue
		}

		if now-price.TSR > thresholdMillis {
			err := cache.Delete(ctx, key)
			if err != nil {
				log.Printf("Ошибка при удалении ключа %s: %v", key, err)
			}
		}
	}
}

func GetAllPrices(t time.Time, cache out.CacheClient) {
	ctx := context.Background()
	var allKeys []string
	var cursor uint64
	maxIterations := 1000
	iteration := 0

	for {
		if iteration >= maxIterations {
			log.Printf("Достигнуто максимальное количество итераций: %d", maxIterations)
			break
		}

		keys, newCursor, err := cache.Scan(ctx, cursor, "*:*", 10000)
		if err != nil {
			log.Printf("Ошибка SCAN: %v", err)
			break
		}

		allKeys = append(allKeys, keys...)
		k := time.Since(t).Seconds()
		fmt.Printf("Итерация %d: cursor=%d, получено ключей: %d, за %f\n", iteration, cursor, len(keys), k)

		cursor = newCursor
		iteration++

		if cursor == 0 {
			break
		}
	}

	fmt.Printf("ИТОГО найдено ключей: %d за %d итераций\n", len(allKeys), iteration)
}
