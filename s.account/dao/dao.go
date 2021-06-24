package dao

import (
	"context"

	"swallowtail/libraries/sql"

	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	db sql.Database
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
