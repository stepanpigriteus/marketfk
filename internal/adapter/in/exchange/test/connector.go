package test

import (
	"context"
	"log"
	"marketfuck/internal/domain/model"
	"math/rand"
	"sync"
	"time"
)

type TestGenerator struct {
	Exchange     model.Exchange
	Pairs        []string
	BasePrice    map[string]float64
	PriceChanges map[string]float64
	Ticker       *time.Ticker
	Mu           sync.RWMutex
}

func NewTestGenerator(exchangeName string) *TestGenerator {
	pairs := []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}

	basePrice := map[string]float64{
		"BTCUSDT":  45000.0,
		"DOGEUSDT": 0.08,
		"TONUSDT":  2.5,
		"SOLUSDT":  100.0,
		"ETHUSDT":  3000.0,
	}

	priceChanges := map[string]float64{
		"BTCUSDT":  0.01,  // 1%
		"DOGEUSDT": 0.05,  // 5%
		"TONUSDT":  0.03,  // 3%
		"SOLUSDT":  0.02,  // 2%
		"ETHUSDT":  0.015, // 1.5%
	}

	return &TestGenerator{
		Exchange:     model.Exchange{Name: exchangeName},
		Pairs:        pairs,
		BasePrice:    basePrice,
		PriceChanges: priceChanges,
		Ticker:       time.NewTicker(100 * time.Millisecond), // Генерация каждые 100 мс
	}
}

// генерирует синтетические данные и отправляет их в канал
func GenConnectAndRead(exchangeName string, wg *sync.WaitGroup, output chan<- model.Price) {
	defer wg.Done()

	log.Printf("Запуск тестового генератора для %s", exchangeName)

	generator := NewTestGenerator(exchangeName)
	defer generator.Ticker.Stop()

	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	for {
		select {
		case <-ctx.Done():
			log.Printf("Тестовый генератор %s остановлен", exchangeName)
			return
		case <-generator.Ticker.C:
			// Генерируем цены для всех пар
			for _, pair := range generator.Pairs {
				price := generator.generatePrice(pair)

				select {
				case output <- price:
				case <-ctx.Done():
					log.Printf("Контекст отменен для генератора %s", exchangeName)
					return
				default:
					log.Printf("Канал заполнен для %s:%s, пропускаем обновление", exchangeName, pair)
				}
			}
		}
	}
}

func (g *TestGenerator) generatePrice(pair string) model.Price {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	basePrice := g.BasePrice[pair]
	maxChange := g.PriceChanges[pair]

	// Генерируем случайное изменение цены в пределах ±maxChange%
	changePercent := (rand.Float64() - 0.5) * 2 * maxChange
	priceChange := basePrice * changePercent

	newPrice := basePrice + priceChange

	// Обновляем базовую цену для следующей итерации (эмуляция трендов)
	g.BasePrice[pair] = newPrice

	// Добавляем некоторую волатильность
	volatility := rand.Float64() * 0.001 // 0.1% случайной волатильности
	if rand.Float64() > 0.5 {
		newPrice += newPrice * volatility
	} else {
		newPrice -= newPrice * volatility
	}

	// Убеждаемся, что цена положительная
	if newPrice <= 0 {
		newPrice = basePrice * 0.9 // Минимум 90% от базовой цены
	}

	return model.Price{
		PairName:  pair,
		Price:     newPrice,
		Exchange:  g.Exchange.Name,
		Timestamp: int64(time.Now().Second()),
	}
}

// запускает несколько тестовых генераторов
func StartTestGenerators(wg *sync.WaitGroup, output chan<- model.Price) {
	exchanges := []string{"TestExchange1", "TestExchange2", "TestExchange3"}

	for _, exchangeName := range exchanges {
		wg.Add(1)
		go GenConnectAndRead(exchangeName, wg, output)
	}
}

// запускает один тестовый генератор с конфигурируемыми параметрами
func StartSingleTestGenerator(exchangeName string, pairs []string, interval time.Duration, wg *sync.WaitGroup, output chan<- model.Price) {
	defer wg.Done()

	log.Printf("Запуск конфигурируемого тестового генератора для %s", exchangeName)

	if len(pairs) == 0 {
		pairs = []string{"BTCUSDT", "DOGEUSDT", "TONUSDT", "SOLUSDT", "ETHUSDT"}
	}

	generator := NewTestGenerator(exchangeName)
	generator.Pairs = pairs
	generator.Ticker.Stop() // Останавливаем старый тикер
	generator.Ticker = time.NewTicker(interval)
	defer generator.Ticker.Stop()

	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	for {
		select {
		case <-ctx.Done():
			log.Printf("Конфигурируемый тестовый генератор %s остановлен", exchangeName)
			return
		case <-generator.Ticker.C:
			for _, pair := range generator.Pairs {
				price := generator.generatePrice(pair)

				select {
				case output <- price:
				case <-ctx.Done():
					return
				default:
					// Канал заполнен, можно добавить логику буферизации или просто пропустить
				}
			}
		}
	}
}

// генерирует исторические данные для инициализации системы
func GenerateHistoricalData(exchangeName string, pairs []string, duration time.Duration, interval time.Duration) []model.Price {
	generator := NewTestGenerator(exchangeName)
	generator.Pairs = pairs

	var prices []model.Price
	startTime := time.Now().Add(-duration)

	for t := startTime; t.Before(time.Now()); t = t.Add(interval) {
		for _, pair := range pairs {
			price := generator.generatePrice(pair)
			price.Timestamp = t
			prices = append(prices, price)
		}
	}

	log.Printf("Сгенерировано %d исторических записей для %s", len(prices), exchangeName)
	return prices
}
