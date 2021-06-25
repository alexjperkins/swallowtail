package sql

import (
	"context"
	"fmt"
	"strconv"
	"swallowtail/libraries/util"
	"sync"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/opentracing/opentracing-go"
)

var (
	connUrl string
	mu      sync.Mutex
)

func init() {
	connUrl = util.SetEnv("SWALLOWTAIL_POSTGRES_CONNECTION_URL")
}

// NewPostgresSQL creates a new postgres database connection.
// TODO: find a way to not pass the service name but to find it dynamically; maybe from hostname?
func NewPostgresSQL(ctx context.Context, applySchema bool, serviceName string) (Database, error) {
	cfg, err := pgconn.ParseConfig(connUrl)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to parse config from postgres connection URL.", map[string]string{
			"postgres_connection_url": connUrl,
		})
	}

	errParams := map[string]string{
		"user": cfg.User,
		"host": cfg.Host,
		"port": strconv.Itoa(int(cfg.Port)),
		"db":   cfg.Database,
	}

	url := fmt.Sprintf("postgresql://%s:%s@%s:%v/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to establish connection to postgresql database", errParams)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, terrors.Augment(err, "Failed to reach to database: opts", errParams)
	}

	slog.Debug(ctx, "Established connection to postgres database", errParams)

	if applySchema {
		err = CreateSchema(ctx, pool, serviceName)
		if err != nil {
			pool.Close()
			return nil, err
		}
	}

	// Background closer, close on cancelled context.
	go func() {
		select {
		case <-ctx.Done():
			slog.Debug(context.TODO(), "Closing postgres connection", map[string]string{
				"service_name": serviceName,
			})
			pool.Close()
		}
	}()

	return &psql{p: pool}, nil
}

type psql struct {
	p *pgxpool.Pool
}

// TODO: add metrics.
func (p *psql) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Postgres query.")
	defer span.Finish()
	return p.p.Query(ctx, sql, args...)
}

func (p *psql) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Postgres execution.")
	defer span.Finish()
	return p.p.Exec(ctx, sql, args...)
}

func (p *psql) Select(ctx context.Context, dest interface{}, sql string, args ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Postgres select statement.")
	defer span.Finish()
	return pgxscan.Select(ctx, p.p, dest, sql, args...)
}

func (p *psql) Transaction(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Postgres transaction.")
	defer span.Finish()
	return p.p.BeginTx(ctx, txOptions)
}
