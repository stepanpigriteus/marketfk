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

func (s *priceService) GetHighestPrice(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	return price, nil
}

func (s *priceService) GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetLowestPrice(ctx context.Context, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.Price, error) {
	var price model.Price
	return price, nil
}

func (s *priceService) GetAveragePrice(ctx context.Context, pairName string, period time.Duration) (float64, error) {
	var avg float64
	return avg, nil
}

func (s *priceService) GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string, period time.Duration) (float64, error) {
	var avg float64
	return avg, nil
}

// это допик для кеша
func (s *priceService) GetHighestPriceFromCache(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var result model.AggregatedPrice

	periodMs := int64(period.Milliseconds())
	redisData, _ := GetAllPrices(s.redisClient, periodMs)

	if len(redisData) == 0 {
		return model.AggregatedPrice{}, fmt.Errorf("данные из кеша отсутствуют")
	}

	cutoff := time.Now().Add(-period)
	result.MaxPrice = -1

	for _, price := range redisData {
		fmt.Println(price.Timestamp)
		if price.PairName == pairName && !price.Timestamp.Before(cutoff) {
			if price.MaxPrice > result.MaxPrice {
				result = price
			}
		}
	}

	if result.MaxPrice == -1 || result.Timestamp.Before(cutoff) {
		return model.AggregatedPrice{}, fmt.Errorf("не найдено цен в кеше за указанный период")
	}

	return result, nil
}
