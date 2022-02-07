package sql

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
	"github.com/opentracing/opentracing-go"

	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/util"
)

const (
	maxConnectionAttempts = 5
)

var (
	connURL string
	mu      sync.Mutex
)

func init() {
	connURL = util.SetEnv("SWALLOWTAIL_POSTGRES_CONNECTION_URL")
}

// NewPostgresSQL creates a new postgres database connection.
// TODO: find a way to not pass the service name but to find it dynamically; maybe from hostname?
func NewPostgresSQL(ctx context.Context, applySchema bool, serviceName string) (Database, error) {
	cfg, err := pgconn.ParseConfig(connURL)
	if err != nil {
		return nil, terrors.Augment(err, "Failed to parse config from postgres connection URL.", map[string]string{
			"postgres_connection_url": connURL,
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

	// Attempt connection with retry.
	if err := pingWithRetry(ctx, pool); err != nil {
		return nil, gerrors.Augment(err, "failed_to_establish_connection_to_postgres", errParams)
	}

	slog.Debug(ctx, "Established connection to postgres database", errParams)

	// Apply schema.
	if applySchema {
		if err = CreateSchema(ctx, pool, serviceName); err != nil {
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

func (p *psql) Get(ctx context.Context, destination interface{}, query string, args ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Postgres get statement")
	defer span.Finish()
	return p.p.QueryRow(ctx, query, args...).Scan(destination)
}

func (p *psql) Transaction(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Postgres transaction.")
	defer span.Finish()
	return p.p.BeginTx(ctx, txOptions)
}

func pingWithRetry(ctx context.Context, pool *pgxpool.Pool) error {
	var (
		cErr error
		boff = backoff.NewExponentialBackOff()
	)
	for i := 0; i < maxConnectionAttempts; i++ {
		if err := pool.Ping(ctx); err != nil {
			cErr = multierror.Append(err)
			d := boff.NextBackOff()
			slog.Trace(ctx, "Failed to connect to postgres instance, retrying...", map[string]string{
				"attempt":        strconv.Itoa(i),
				"sleep_duration": d.String(),
			})
			time.Sleep(d)
		}
	}
	if cErr != nil {
		return gerrors.Augment(cErr, "failed_to_establish_connection_to_postgres_instance.after_retries", nil)
	}

	return nil
}
