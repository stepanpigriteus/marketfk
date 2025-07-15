package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"marketfuck/internal/adapter/out_impl_for_port_out/cache/redis"
	"marketfuck/internal/domain/model"
	"sync"
	"sync/atomic"
	"time"
)

type Mode int

const (
	LiveMode Mode = iota
	TestMode
)

type ModeServiceImpl struct {
	mu            sync.Mutex
	currentMode   Mode
	cancelFunc    context.CancelFunc
	wg            *sync.WaitGroup
	redis         *redis.RedisCache
	counter       *atomic.Uint64
	outputChan    <-chan model.Price
	aggregatorFn  func(ctx context.Context, counter *atomic.Uint64, redis redis.RedisCache, ports []string) <-chan model.Price
	fakeExchanges map[string]context.CancelFunc
}

// Конструктор
func NewModeService(redis *redis.RedisCache, counter *atomic.Uint64, aggregatorFn func(context.Context, *atomic.Uint64, redis.RedisCache, []string) <-chan model.Price) *ModeServiceImpl {
	return &ModeServiceImpl{
		redis:         redis,
		counter:       counter,
		aggregatorFn:  aggregatorFn,
		wg:            &sync.WaitGroup{},
		fakeExchanges: make(map[string]context.CancelFunc),
	}
}

func (m *ModeServiceImpl) SwitchToTestMode(ctx context.Context) error {
	return m.switchMode(ctx, TestMode, []string{"50201", "50202", "50203"})
}

func (m *ModeServiceImpl) SwitchToLiveMode(ctx context.Context) error {
	return m.switchMode(ctx, LiveMode, []string{"40101", "40102", "40103"})
}

func (m *ModeServiceImpl) switchMode(ctx context.Context, newMode Mode, ports []string) error {
	log.Println("switchMode() called with mode:", newMode)
	m.mu.Lock()
	defer m.mu.Unlock()

	// Остановка предыдущей агрегации, даже если режим тот же
	if m.cancelFunc != nil {
		log.Println("Stopping previous aggregation")
		m.cancelFunc()
		m.wg.Wait()
		log.Println("Предыдущая агрегация остановлена")
	}

	// Очистка Redis
	log.Println("Clearing Redis cache")
	if err := m.redis.Clear(ctx); err != nil {
		log.Printf("Ошибка очистки Redis: %v", err)
		return fmt.Errorf("ошибка очистки Redis: %w", err)
	}

	// Новый контекст и запуск новой агрегации
	log.Println("Creating new context for aggregation")
	newCtx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel
	m.currentMode = newMode

	// Запуск агрегации
	log.Println("Starting aggregatorFn with ports:", ports)
	m.outputChan = m.aggregatorFn(newCtx, m.counter, *m.redis, ports)
	if m.outputChan == nil {
		log.Println("! aggregatorFn вернул nil канал")
		return errors.New("aggregatorFn returned nil")
	}
	log.Println("aggregatorFn returned valid channel")

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		log.Println("Started reading from outputChan")
		for price := range m.outputChan {
			price.TSR = time.Now().UnixMilli()
			key := fmt.Sprintf("%s:%s:%d", price.PairName, price.Exchange, price.TSR)
			if err := m.redis.SetPrice(newCtx, key, price, 0); err != nil {
				log.Printf("Ошибка установки цены: %v", err)
			}
		}
		log.Println("outputChan closed")
	}()

	log.Printf("Переключено на режим: %v", newMode)
	return nil
}
