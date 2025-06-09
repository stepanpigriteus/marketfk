package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"marketfuck/internal/domain/model"
	"strings"
	"time"
)

type PriceRepository struct {
	db *sql.DB
}

func NewPriceRepository(db *sql.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

// написать нормальную функцию вставки данных!!!
func (r *PriceRepository) SavePrice(ctx context.Context, prices []model.AggregatedPrice) error {
	if len(prices) == 0 {
		return nil
	}

	query := `INSERT INTO aggregated_prices (pair_name, exchange, timestamp, average_price, min_price, max_price) VALUES `
	var args []interface{}
	var valueStrings []string

	for i, p := range prices {
		base := i * 6
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			base+1, base+2, base+3, base+4, base+5, base+6,
		))

		args = append(args,
			p.PairName,
			p.Exchange,
			p.Timestamp,
			p.AveragePrice,
			p.MinPrice,
			p.MaxPrice,
		)
	}

	query += strings.Join(valueStrings, ", ")

	// 🔧 Без троеточия — передаём args как []interface{}
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PriceRepository) GetLatestPrice(ctx context.Context, pairName string) (model.Price, error) {
	return model.Price{}, nil
}

func (r *PriceRepository) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.Price, error) {
	return model.Price{}, nil
}

func (r *PriceRepository) GetPricesInPeriod(ctx context.Context, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
	return nil, nil
}

func (r *PriceRepository) GetPricesInPeriodByExchange(ctx context.Context, exchangeID, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
	return nil, nil
}
