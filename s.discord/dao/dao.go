package dao

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	db *pgxpool.Pool
	mu sync.RWMutex
)

// Init creates the database connection.
func Init(ctx context.Context, opts *pgconn.Config) (func(), error) {
	// Connect to the database
	slogParams := map[string]string{
		"user": opts.User,
		"host": opts.Host,
		"port": strconv.Itoa(int(opts.Port)),
		"db":   opts.Database,
	}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%v/%s", opts.User, opts.Password, opts.Host, opts.Port, opts.Database)
	dbpool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to connect to database: opts", slogParams)
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, terrors.Augment(err, "Failed to reach to database: opts", slogParams)
	}

	slog.Info(ctx, "Established connection to database", slogParams)

	// Create the initial schema for the database.
	err = createSchema(ctx, dbpool)
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	db = dbpool
	return dbpool.Close, nil
}

func DB() *pgxpool.Pool {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

func createSchema(ctx context.Context, dbpool *pgxpool.Pool) error {
	sql, err := loadSQLFile(ctx)
	if err != nil {
		return err
	}

	_, err = dbpool.Exec(ctx, sql)
	if err != nil {
		return terrors.Augment(err, "Failed to create initial postgres schema", nil)
	}

	return nil
}

func loadSQLFile(ctx context.Context) (string, error) {
	c, err := ioutil.ReadFile("/home/alexjperkins/repos/swallowtail/s.discord/config/postgres.sql")
	if err != nil {
		return "", terrors.Augment(err, "Failed to load config postgres sql file.", nil)
	}
	slog.Info(ctx, "Loaded config postgres sql file.")
	return string(c), nil
}
