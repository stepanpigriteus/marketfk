package live

import (
	"bufio"
	"context"
	"encoding/json"
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

	subscriptions map[string][]chan<- model.Price
	subsMu        sync.RWMutex

	latestPrices map[string]model.Price
	pricesMu     sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	reconnectCh chan struct{}
}

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

func (e *LiveExchangeClient) StartReading(output chan<- model.Price) {
	go func() {
		defer e.conn.Close()

		scanner := bufio.NewScanner(e.conn)
		for scanner.Scan() {
			var price model.Price
			err := json.Unmarshal(scanner.Bytes(), &price)
			if err != nil {
				continue // пропускаем ошибку
			}

			price.Exchange = e.exchange.Name

			output <- price
		}
	}()
}
