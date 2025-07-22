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
		SELECT pair_name, exchange, max_price, timestamp
		FROM aggregated_prices
		WHERE pair_name = $1
		ORDER BY average_price DESC
		LIMIT 1;
	`
	err := r.db.QueryRowContext(ctx, query, pairName).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MaxPrice,
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

func (r *PriceRepository) GetHighestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	query := `
		SELECT pair_name, exchange, max_price, timestamp
		FROM aggregated_prices
		WHERE exchange = $1 AND pair_name = $2
		ORDER BY max_price DESC
		LIMIT 1;
	`
	err := r.db.QueryRowContext(ctx, query, exchangeID, pairName).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MaxPrice,
		&price.Timestamp,
	)

	fmt.Println(err, price)
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

func (r *PriceRepository) GetHighestPriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	lastTimestamp, err := r.GetLastRecordTime(ctx, pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	endTime := lastTimestamp
	startTime := endTime.Add(-period)

	query := `
	SELECT pair_name, exchange, MAX(max_price) as max_price, timestamp
	FROM aggregated_prices
	WHERE pair_name = $1 AND timestamp >= $2 AND timestamp <= $3 AND exchange =$4
	GROUP BY pair_name, exchange, timestamp
	ORDER BY max_price DESC
	LIMIT 1;
	`
	err = r.db.QueryRowContext(ctx, query, pairName, startTime, endTime, exchangeID).Scan(
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

func (r *PriceRepository) GetLovestPrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	query := `
		SELECT pair_name, exchange, min_price, timestamp
		FROM aggregated_prices
		WHERE pair_name = $1
		ORDER BY average_price ASC
		LIMIT 1;
	`
	err := r.db.QueryRowContext(ctx, query, pairName).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MinPrice,
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

func (s *PriceRepository) GetLowestPriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	lastTimestamp, err := s.GetLastRecordTime(ctx, pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	endTime := lastTimestamp
	startTime := endTime.Add(-period)
	query := `
	SELECT pair_name, exchange, MIN(min_price) as min_price, timestamp
	FROM aggregated_prices
	WHERE pair_name = $1 AND timestamp >= $2 AND timestamp <= $3
	GROUP BY pair_name, exchange, timestamp
	ORDER BY min_price ASC
	LIMIT 1;
`

	err = s.db.QueryRowContext(ctx, query, pairName, startTime, endTime).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MinPrice,
		&price.Timestamp,
	)
	if err != nil {
		return model.AggregatedPrice{}, fmt.Errorf("could not retrieve highest price: %v", err)
	}
	fmt.Println(price)
	return price, nil
}

func (r *PriceRepository) GetLowestPriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	query := `
		SELECT pair_name, exchange, min_price, timestamp
		FROM aggregated_prices
		WHERE exchange = $1 AND pair_name = $2
		ORDER BY min_price ASC
		LIMIT 1;
	`
	err := r.db.QueryRowContext(ctx, query, exchangeID, pairName).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MinPrice,
		&price.Timestamp,
	)

	fmt.Println(err, price)
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

func (r *PriceRepository) GetLowestPriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	lastTimestamp, err := r.GetLastRecordTime(ctx, pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}
	endTime := lastTimestamp
	startTime := endTime.Add(-period)

	query := `
	SELECT pair_name, exchange, MIN(min_price) as min_price, timestamp
	FROM aggregated_prices
	WHERE pair_name = $1 AND timestamp >= $2 AND timestamp <= $3 AND exchange =$4
	GROUP BY pair_name, exchange, timestamp
	ORDER BY min_price ASC
	LIMIT 1;
	`
	err = r.db.QueryRowContext(ctx, query, pairName, startTime, endTime, exchangeID).Scan(
		&price.PairName,
		&price.Exchange,
		&price.MinPrice,
		&price.Timestamp,
	)
	if err != nil {
		return model.AggregatedPrice{}, fmt.Errorf("could not retrieve highest price: %v", err)
	}
	fmt.Println(price)
	return price, nil
}

func (r *PriceRepository) GetAveragePrice(ctx context.Context, pairName string) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice

	query := `
		SELECT AVG(average_price)
		FROM aggregated_prices
		WHERE pair_name = $1;
	`

	err := r.db.QueryRowContext(ctx, query, pairName).Scan(&price.AveragePrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AggregatedPrice{}, fmt.Errorf("no data found for pair %s", pairName)
		}
		return model.AggregatedPrice{}, fmt.Errorf("failed to compute average: %v", err)
	}

	price.PairName = pairName
	price.Timestamp = time.Now() // или возьми последний timestamp, если нужно

	fmt.Printf("Computed average price: PairName=%s, AveragePrice=%f\n",
		price.PairName, price.AveragePrice)

	return price, nil
}

func (s *PriceRepository) GetAveragePriceInPeriod(ctx context.Context, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice
	lastTimestamp, err := s.GetLastRecordTime(ctx, pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	endTime := lastTimestamp
	startTime := endTime.Add(-period)

	query := `
	SELECT AVG(average_price)
	FROM aggregated_prices
	WHERE pair_name = $1 AND timestamp BETWEEN $2 AND $3;
	`

	err = s.db.QueryRowContext(ctx, query, pairName, startTime, endTime).Scan(&price.AveragePrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AggregatedPrice{}, fmt.Errorf("no data found for pair %s", pairName)
		}
		return model.AggregatedPrice{}, fmt.Errorf("failed to compute average: %v", err)
	}

	price.PairName = pairName
	price.Timestamp = endTime

	fmt.Printf("Average price for %s in last %s = %f\n", pairName, period.String(), price.AveragePrice)

	return price, nil
}

func (r *PriceRepository) GetAveragePriceByExchange(ctx context.Context, exchangeID, pairName string) (model.AggregatedPrice, error) {
	query := `
		SELECT average_price
		FROM aggregated_prices
		WHERE exchange = $1 AND pair_name = $2;
	`

	rows, err := r.db.QueryContext(ctx, query, exchangeID, pairName)
	if err != nil {
		return model.AggregatedPrice{}, fmt.Errorf("ошибка запроса к БД: %v", err)
	}
	defer rows.Close()

	var sum float64
	var count int

	for rows.Next() {
		var avg float64
		if err := rows.Scan(&avg); err != nil {
			return model.AggregatedPrice{}, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		sum += avg
		count++
	}

	if err := rows.Err(); err != nil {
		return model.AggregatedPrice{}, fmt.Errorf("ошибка обхода строк: %v", err)
	}

	if count == 0 {
		return model.AggregatedPrice{}, fmt.Errorf("нет данных для пары %s на бирже %s", pairName, exchangeID)
	}

	average := sum / float64(count)

	return model.AggregatedPrice{
		PairName:     pairName,
		Exchange:     exchangeID,
		AveragePrice: average,
	}, nil
}

func (r *PriceRepository) GetAveragePriceByExchangeInPeriod(ctx context.Context, exchangeID, pairName string, period time.Duration) (model.AggregatedPrice, error) {
	var price model.AggregatedPrice

	lastTimestamp, err := r.GetLastRecordTime(ctx, pairName)
	if err != nil {
		return model.AggregatedPrice{}, err
	}

	endTime := lastTimestamp
	startTime := endTime.Add(-period)

	query := `
		SELECT AVG(average_price)
		FROM aggregated_prices
		WHERE pair_name = $1 AND exchange = $2 AND timestamp BETWEEN $3 AND $4;
	`

	err = r.db.QueryRowContext(ctx, query, pairName, exchangeID, startTime, endTime).Scan(&price.AveragePrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AggregatedPrice{}, fmt.Errorf("нет данных для %s на %s", pairName, exchangeID)
		}
		return model.AggregatedPrice{}, fmt.Errorf("ошибка выборки: %v", err)
	}

	price.PairName = pairName
	price.Exchange = exchangeID
	price.Timestamp = endTime

	fmt.Printf("Средняя цена %s на %s за период %s: %f\n", pairName, exchangeID, period.String(), price.AveragePrice)
	return price, nil
}
