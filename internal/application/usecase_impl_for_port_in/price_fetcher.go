package usecase

import (
	"context"
	"fmt"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/adapter/out_impl_for_port_out/storage/postgres"
	"marketfuck/internal/application/port/in"
	"marketfuck/internal/domain/model"
	"marketfuck/pkg/utils"
	"time"
)

type priceService struct {
	priceRepo   postgres.PriceRepository
	redisClient *redis.RedisCache
}

func NewPriceService(repo postgres.PriceRepository, redisClient *redis.RedisCache) in.PriceService {
	return &priceService{priceRepo: repo, redisClient: redisClient}
}

func (s *priceService) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	price, err := s.priceRepo.GetLatestPrice(ctx, name)
	if err != nil {
		fmt.Printf("Error in GetLatestPrice: %v\n", err) // Логируем ошибку
		return model.AggregatedPrice{}, err
	}

	// Возвращаем результат
	return price, nil
}

func (s *priceService) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	price, err := s.priceRepo.GetLatestPriceByExchange(ctx, exchangeID, name)
	if err != nil {
		fmt.Printf("Error in GetLatestPrice: %v\n", err)
		return model.AggregatedPrice{}, err
	}

	return price, nil
}

func (s *priceService) GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	price, err = s.priceRepo.GetHighestPriceInPeriod(ctx, name, period)
	if err != nil {
		fmt.Printf("Error in GetHighestPriceInPeriod: %v\n", err)
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetHighestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetHighestPrice(ctx, name)
	if err != nil {
		fmt.Printf("Error in GetHighestPrice: %v\n", err)
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetHighestPriceByExchange(ctx, exchangeID, name)
	return price, nil
}

func (s *priceService) GetHighestPriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetHighestPriceByExchangeInPeriod(ctx, exchangeID, name, period)
	return price, nil
}

// это допик для кеша
func (s *priceService) GetHighestPriceFromCache(ctx context.Context, pairName string, period time.Duration, exchange string, flag string) (model.AggregatedPrice, error) {
	var result model.AggregatedPrice
	periodMs := int64(period.Milliseconds())
	redisData, _ := GetAllPrices(s.redisClient, periodMs)
	if len(redisData) == 0 {
		return model.AggregatedPrice{}, fmt.Errorf("данные из кеша отсутствуют")
	}

	var latestTime time.Time
	totalPairRecords := 0
	fmt.Println(redisData)
	for _, price := range redisData {
		if price.PairName == pairName {
			totalPairRecords++
			if price.Timestamp.After(latestTime) {
				latestTime = price.Timestamp
			}
		}
	}

	if latestTime.IsZero() {
		return model.AggregatedPrice{}, fmt.Errorf("нет данных для пары %s", pairName)
	}

	cutoff := latestTime.Add(-period)

	validCount := 0
	var sum float64
	var avgCount int
	var initialized bool

	for _, price := range redisData {
		if price.PairName == pairName && (exchange == "" || price.Exchange == exchange) {
			if !price.Timestamp.Before(cutoff) && !price.Timestamp.After(latestTime) {
				validCount++

				switch flag {
				case "max":
					if !initialized || price.MaxPrice > result.MaxPrice {
						result = price
						initialized = true
					}
				case "min":
					if !initialized || price.MinPrice < result.MinPrice {
						result = price
						initialized = true
					}

				case "avg":
					sum += price.AveragePrice
					avgCount++
				}
			}
		}
	}

	fmt.Printf("Found %d valid records for %s\n", validCount, pairName)

	if flag == "avg" {
		if avgCount == 0 {
			return model.AggregatedPrice{}, fmt.Errorf("не найдено цен в кеше за указанный период")
		}
		result = model.AggregatedPrice{
			PairName:  pairName,
			Exchange:  exchange,
			MaxPrice:  sum / float64(avgCount),
			Timestamp: latestTime,
		}
	} else if !initialized {
		return model.AggregatedPrice{}, fmt.Errorf("не найдено цен в кеше за указанный период")
	}

	fmt.Printf("Best record: timestamp=%v, price=%v\n", result.Timestamp, result.MaxPrice)
	return result, nil
}

func (s *priceService) GetLowestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetLovestPrice(ctx, name)
	if err != nil {
		fmt.Printf("Error in GetHighestPrice: %v\n", err)
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetLowestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	price, err = s.priceRepo.GetLowestPriceInPeriod(ctx, name, period)
	if err != nil {
		fmt.Printf("Error in GetHighestPriceInPeriod: %v\n", err)
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetLowestPriceByExchange(ctx, exchangeID, name)
	return price, nil
}

func (s *priceService) GetLowestPriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetLowestPriceByExchangeInPeriod(ctx, exchangeID, name, period)
	return price, nil
}

func (s *priceService) GetAveragePrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetAveragePrice(ctx, name)
	if err != nil {
		fmt.Printf("Error in GetHighestPrice: %v\n", err)
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetAveragePriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	price, err = s.priceRepo.GetAveragePriceInPeriod(ctx, name, period)
	if err != nil {
		fmt.Printf("Error in GetHighestPriceInPeriod: %v\n", err)
		return model.AggregatedPrice{}, err
	}
	return price, nil
}

func (s *priceService) GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetAveragePriceByExchange(ctx, exchangeID, name)
	return price, nil
}

func (s *priceService) GetAveragePriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	name, err := utils.PairNameValidFormatter(pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	price, err = s.priceRepo.GetAveragePriceByExchangeInPeriod(ctx, exchangeID, name, period)
	return price, nil
}
