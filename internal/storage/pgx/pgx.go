package pgx

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"lead-bitrix/internal/config"
	"time"
)

//go:embed sql/tables.sql
var tableSchema string

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(cfg *config.Config) (*Storage, error) {

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBConfig.Username,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Name,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		err = fmt.Errorf("error parsing connection string %s", err.Error())
		return nil, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	if err := createTable(ctx, db); err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	return &Storage{db: db}, nil

}

func createTable(ctx context.Context, db *pgxpool.Pool) error {

	commandTag, err := db.Exec(ctx, tableSchema)

	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	fmt.Println(commandTag.RowsAffected())
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
