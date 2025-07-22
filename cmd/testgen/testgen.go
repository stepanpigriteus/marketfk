package testgen

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"marketfuck/internal/domain/model"
	"math/rand"
	"net"
	"time"
)

var pairs = []string{
	"BTCUSDT",
	"DOGEUSDT",
	"TONUSDT",
	"SOLUSDT",
	"ETHUSDT",
}

func RandomPrice(pair string) float64 {
	switch pair {
	case "BTCUSDT":
		return 29000 + rand.Float64()*2000
	case "ETHUSDT":
		return 1800 + rand.Float64()*300
	case "DOGEUSDT":
		return 0.06 + rand.Float64()*0.02
	case "TONUSDT":
		return 1.5 + rand.Float64()*0.5
	case "SOLUSDT":
		return 20 + rand.Float64()*10
	default:
		return 100.0 + rand.Float64()*50
	}
}

func GeneratePrice(exchangeName string) model.Price {
	pair := pairs[rand.Intn(len(pairs))]
	ts := time.Now().UnixMilli()
	return model.Price{
		PairName:  pair,
		Exchange:  exchangeName,
		Price:     RandomPrice(pair),
		Timestamp: ts,
		TSR:       ts,
	}
}

func StartFakeExchangeWithCtx(ctx context.Context, port string, exchangeName string) error {
	addr := "127.0.0.1:" + port

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("не удалось запустить сервер на порту %s: %w", port, err)
	}
	log.Printf("[%s] слушает на %s", exchangeName, addr)

	go func() {
		<-ctx.Done()
		log.Printf("[%s] остановка сервера на порту %s", exchangeName, port)
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil // корректное завершение
			default:
				log.Printf("[%s] ошибка соединения: %v", exchangeName, err)
				continue
			}
		}
		go handleExchangeConnection(ctx, conn, exchangeName)
	}
}

func handleExchangeConnection(ctx context.Context, conn net.Conn, exchangeName string) {
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] завершение соединения", exchangeName)
			return
		case <-ticker.C:
			price := GeneratePrice(exchangeName)
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := encoder.Encode(price); err != nil {
				log.Printf("[%s] ошибка отправки данных: %v", exchangeName, err)
				return
			}

			if err := encoder.Encode(price); err != nil {
				log.Printf("[%s] ошибка отправки данных: %v", exchangeName, err)
				return
			}
		}
	}
}
