package live

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
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
		defer func() {
			e.conn.Close()
			close(output)
		}()

		scanner := bufio.NewScanner(e.conn)
		for scanner.Scan() {
			line := scanner.Bytes()

			var price model.Price
			err := json.Unmarshal(line, &price)
			if err != nil {
				log.Printf("JSON unmarshal error: %v", err)
				continue
			}
			price.Exchange = e.exchange.Name
			output <- price
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Scanner error: %v", err)
		} else {
			log.Printf("Scanner finished without error")
		}
	}()
}

func (e *LiveExchangeClient) Connect() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.isConnected {
		return nil
	}

	conn, err := net.Dial("tcp", e.address)
	if err != nil {
		return err
	}

	e.conn = conn
	e.isConnected = true
	return nil
}
