package in

import (
	"context"
	"marketfuck/internal/domain/model"
)


type ExchangeClient interface {
	//получает информацию о бирже
	GetExchangeInfo(ctx context.Context) (model.Exchange, error)
	//получает последнюю цену для указанной пары
	GetLatestPrice(ctx context.Context, pairName string) (model.Price, error)
	// подписывается на обновления цен для указанной пары
	SubscribePriceUpdates(ctx context.Context, pairName string, ch chan<- model.Price) error
	// отписывается от обновлений цен
	UnsubscribePriceUpdates(ctx context.Context, pairName string) error
	// проверяет, установлено ли соединение с биржей
	CheckConnection(ctx context.Context) (bool, error)
}