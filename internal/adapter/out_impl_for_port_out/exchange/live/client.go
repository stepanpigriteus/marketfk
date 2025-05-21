package live

import (
	"context"
	"marketfuck/internal/domain/model"
	"net"
	"sync"
)

type LiveExchangeClient struct {
	exchange    model.Exchange
	address     string
	conn        net.Conn
	isConnected bool
	mu          sync.RWMutex

	// Подписки на обновления цен
	subscriptions map[string][]chan<- model.Price
	subsMu        sync.RWMutex

	// Последние полученные цены для каждой пары
	latestPrices map[string]model.Price
	pricesMu     sync.RWMutex

	// Контекст и канал для остановки
	ctx    context.Context
	cancel context.CancelFunc

	// Канал для переподключения
	reconnectCh chan struct{}
}

// NewLiveExchangeClient создает новый клиент для подключения к бирже
func NewLiveExchangeClient(exchange model.Exchange, address string) *LiveExchangeClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &LiveExchangeClient{
		exchange:      exchange,
		address:       address,
		subscriptions: make(map[string][]chan<- model.Price),
		latestPrices:  make(map[string]model.Price),
		ctx:           ctx,
		cancel:        cancel,
		reconnectCh:   make(chan struct{}, 1),
	}
}

func (e *LiveExchangeClient) GetLatestPrice(ctx context.Context, pairName string) (model.Price, error) {
	return model.Price{}, nil
}

func (e *LiveExchangeClient) GetExchangeInfo(ctx context.Context) (model.Exchange, error) {
	return e.exchange, nil
}

func (e *LiveExchangeClient) SubscribePriceUpdates(ctx context.Context, pairName string, ch chan<- model.Price) error {
	return nil
}

func (e *LiveExchangeClient) UnsubscribePriceUpdates(ctx context.Context, pairName string) error {
	return nil
}

func (e *LiveExchangeClient) CheckConnection(ctx context.Context) (bool, error) {
	return false, nil
}
