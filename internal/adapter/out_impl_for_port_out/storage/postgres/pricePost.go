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

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// протестить
func (r *PriceRepository) GetLatestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var aggregatedPrice model.AggregatedPrice
	query := `
    SELECT pair_name, exchange, average_price, timestamp
    FROM aggregated_prices
    WHERE pair_name = $1 
    ORDER BY timestamp DESC
    LIMIT 1;
    `

	err := r.db.QueryRowContext(ctx, query, pairName).Scan(
		&aggregatedPrice.PairName,
		&aggregatedPrice.Exchange,
		&aggregatedPrice.AveragePrice,
		&aggregatedPrice.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AggregatedPrice{}, fmt.Errorf("no data found for pair %s", pairName)
		}
		return model.AggregatedPrice{}, fmt.Errorf("failed to fetch latest price: %v", err)
	}

	fmt.Printf("Fetched aggregated price: PairName=%s, Exchange=%s, AveragePrice=%f, Timestamp=%s\n",
		aggregatedPrice.PairName, aggregatedPrice.Exchange, aggregatedPrice.AveragePrice, aggregatedPrice.Timestamp)

	return aggregatedPrice, nil
}

func (s *PriceRepository) GetLatestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	query := `
   	SELECT pair_name, exchange, average_price, timestamp
    FROM aggregated_prices
    WHERE exchange = $1 AND pair_name = $2
    ORDER BY timestamp DESC
    LIMIT 1;
    `

	err := s.db.QueryRowContext(ctx, query, exchangeID, pairName).Scan(
		&price.PairName,
		&price.Exchange,
		&price.AveragePrice,
		&price.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AggregatedPrice{}, fmt.Errorf("no data found for exchange %s and pair %s", exchangeID, pairName)
		}
		return model.AggregatedPrice{}, fmt.Errorf("failed to fetch latest price: %v", err)
	}

	return price, nil
}

func (r *PriceRepository) GetHighestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	query := `
		SELECT pair_name, exchange, average_price, timestamp
		FROM aggregated_prices
		WHERE pair_name = $1
		ORDER BY average_price DESC
		LIMIT 1;
	`
	err := r.db.QueryRowContext(ctx, query, pairName).Scan(
		&price.PairName,
		&price.Exchange,
		&price.AveragePrice,
		&price.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AggregatedPrice{}, fmt.Errorf("no data found for pair %s", pairName)
		}
		return model.AggregatedPrice{}, fmt.Errorf("failed to fetch latest price: %v", err)
	}

	fmt.Printf("Fetched aggregated price: PairName=%s, Exchange=%s, AveragePrice=%f, Timestamp=%s\n",
		price.PairName, price.Exchange, price.AveragePrice, price.Timestamp)

	return price, nil

}

func (r *PriceRepository) GetPricesInPeriod(ctx context.Context, pairName string, duration time.Time) ([]model.AggregatedPrice, error) {
	return nil, nil
}

func (r *PriceRepository) GetPricesInPeriodByExchange(ctx context.Context, exchangeID, pairName string, startTime, endTime time.Time) ([]model.Price, error) {
	return nil, nil
}

func (s *PriceRepository) GetHighestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	lastTimestamp, err := s.GetLastRecordTime(ctx, pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	endTime := lastTimestamp
	startTime := endTime.Add(-period)
	query := `
	SELECT pair_name, exchange, MAX(max_price) as max_price, timestamp
	FROM aggregated_prices
	WHERE pair_name = $1 AND timestamp >= $2 AND timestamp <= $3
	GROUP BY pair_name, exchange, timestamp
	ORDER BY max_price DESC
	LIMIT 1;
`

	err = s.db.QueryRowContext(ctx, query, pairName, startTime, endTime).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MaxPrice,
		&price.Timestamp,
	)
	if err != nil {
		return model.AggregatedPrice{}, fmt.Errorf("could not retrieve highest price: %v", err)
	}
	fmt.Println(price)
	return price, nil
}

// получение времени последней записи
func (s *PriceRepository) GetLastRecordTime(ctx context.Context, pairName string) (time.Time, error) {
	queryLastRecordTime := `
		SELECT timestamp
		FROM aggregated_prices
		WHERE pair_name = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var lastTimestamp time.Time
	err := s.db.QueryRowContext(ctx, queryLastRecordTime, pairName).Scan(&lastTimestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not retrieve last record time for pair %s: %v", pairName, err)
	}

	return lastTimestamp, nil
}
