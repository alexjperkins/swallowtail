package sql

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Database is the interface for a postgres connection. The implemententaion details should be hidden.
type Database interface {
	// Exec will execute a sql statement on the underlying postgres database & return a command tag object.
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)

	// Query will perform a query on the underlying database & return a list of database rows.
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)

	// Select takes an array of struct representation of a postgresql database row & marshals into
	// it via executing the passed sql statement.
	Select(ctx context.Context, dest interface{}, sql string, args ...interface{}) error

	// Transaction beings a transaction on the underlying database & return a transaction object.
	Transaction(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}
