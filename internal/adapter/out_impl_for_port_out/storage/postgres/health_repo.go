package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresHealthRepository struct {
	db *sql.DB
}

func NewPostgresHealthRepository(db *sql.DB) *PostgresHealthRepository {
	return &PostgresHealthRepository{
		db: db,
	}
}

func (r *PostgresHealthRepository) CheckConnection(ctx context.Context) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.db.PingContext(ctxTimeout); err != nil {
		return false, fmt.Errorf("ping failed: %w", err)
	}

	var result int
	query := "SELECT 1"
	if err := r.db.QueryRowContext(ctxTimeout, query).Scan(&result); err != nil {
		return false, fmt.Errorf("test query failed: %w", err)
	}

	if result != 1 {
		return false, fmt.Errorf("unexpected test query result: %d", result)
	}

	return true, nil
}
