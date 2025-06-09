package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"marketfuck/internal/domain/model"
	"time"
)

type PriceRepository struct {
	db *sql.DB
}

func NewPriceRepository(db *sql.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

// дописать передачу аггрегантов
func (r *PriceRepository) SavePrice(ctx context.Context, prices []model.AggregatedPrice) error {
	if len(prices) == 0 {
		return nil
	}

	query := `
		INSERT INTO aggregated_prices 
		    (pair_name, exchange, timestamp, average_price, min_price, max_price)
		VALUES 
	`
	args := []interface{}{}
	for i, p := range prices {
		idx := i * 6
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6)
		args = append(args, p.PairName, p.Exchange, p.Timestamp, p.AveragePrice, p.MinPrice, p.MaxPrice)
	}
	query = query[:len(query)-1]

	_, err := r.db.ExecContext(ctx, query, args)
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
