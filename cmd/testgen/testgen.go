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
	addr := fmt.Sprintf("marketfuck:%s", port) // Используем имя сервиса в Docker-сети
	for attempt := 1; attempt <= 5; attempt++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			log.Printf("Фейковая биржа %s подключилась к порту %s", exchangeName, port)
			go handleConnection(ctx, conn, exchangeName)
			<-ctx.Done()
			log.Printf("Фейковая биржа %s остановлена", exchangeName)
			return ctx.Err()
		}
		log.Printf("Ошибка подключения фейковой биржи %s к порту %s (попытка %d/5): %v", exchangeName, port, attempt, err)
		select {
		case <-ctx.Done():
			log.Printf("Фейковая биржа %s остановлена по сигналу контекста", exchangeName)
			return ctx.Err()
		case <-time.After(time.Duration(1<<attempt) * time.Second):
			continue
		}
	}
	return fmt.Errorf("не удалось подключиться к порту %s после 5 попыток", port)
}

func handleConnection(ctx context.Context, conn net.Conn, exchangeName string) {
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	log.Printf("New connection established for %s", exchangeName)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Connection closed for %s due to context cancellation", exchangeName)
			return
		case <-ticker.C:
			price := GeneratePrice(exchangeName)
			log.Printf("Sending price from %s: %+v", exchangeName, price)
			if err := encoder.Encode(price); err != nil {
				log.Printf("Ошибка отправки данных с %s: %v", exchangeName, err)
				return
			}
		}
	}
}
