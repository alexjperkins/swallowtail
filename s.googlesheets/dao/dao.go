package dao

import (
	"context"
	"sync"

	"swallowtail/libraries/sql"
	"swallowtail/libraries/sql/mocks"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	db sql.Database
	mu sync.Mutex
)

// Init creates the database connection.
func Init(ctx context.Context, serviceName string) error {
	psql, err := sql.NewPostgresSQL(ctx, true, serviceName)
	if err != nil {
		return terrors.Augment(err, "Failed to initialize dao", map[string]string{
			"service_name": serviceName,
		})
	}
	db = psql
	slog.Debug(ctx, "Dao initialized", map[string]string{
		"service_name": serviceName,
	})
	return nil
}

// WithMock uses a mock db.
func WithMock() {
	if db != nil {
		panic("Cannot set running db as Mock.")
	}
	mu.Lock()
	defer mu.Unlock()

	db = &mocks.Database{}
}
