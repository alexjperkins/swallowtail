package sql

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"

	"swallowtail/libraries/util"
)

const (
	postgresConfigFileName = "postgres.sql"
)

// CreateSchema creates the schema in the local config `postgres.sql` file if it doesn't already exist.
// TODO: diff between postgres & other dbs config files.
func CreateSchema(ctx context.Context, dbpool *pgxpool.Pool, serviceName string) error {
	query, err := loadSQLFile(ctx, serviceName)
	if err != nil {
		return err
	}

	if _, err = dbpool.Exec(ctx, query); err != nil {
		return terrors.Augment(err, "Failed to create initial postgres schema", nil)
	}

	return nil
}

// loadSQLFile loads `<project-root>/<service-name>/config/postgres.sql`.
func loadSQLFile(ctx context.Context, serviceName string) (string, error) {
	root, err := util.RootDir()
	if err != nil {
		return "", terrors.Augment(err, "Failed to load sql config file", map[string]string{
			"service_name": serviceName,
		})
	}

	path := fmt.Sprintf("%s/%s/%s/%s", root, serviceName, "config", postgresConfigFileName)
	slog.Debug(ctx, "Loading postgres sql file", map[string]string{
		"postgres_config_path": path,
	})

	c, err := ioutil.ReadFile(path)
	if err != nil {
		return "", terrors.Augment(err, "Failed to load config postgres sql file.", nil)
	}

	slog.Info(ctx, "Loaded config postgres sql file: %s", path)
	return string(c), nil
}
