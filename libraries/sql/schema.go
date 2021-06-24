package sql

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

const (
	postgresConfigFileName = "postgres.sql"
)

// createSchema creates the schema in the local config `postgres.sql` file if it doesn't already exist.
func createSchema(ctx context.Context, dbpool *pgxpool.Pool, serviceName string) error {
	sql, err := loadSQLFile(ctx, serviceName)
	if err != nil {
		return err
	}

	_, err = dbpool.Exec(ctx, sql)
	if err != nil {
		return terrors.Augment(err, "Failed to create initial postgres schema", nil)
	}

	return nil
}

// loadSQLFile loads `<project-root>/<service-name>/config/postgres.sql`.
func loadSQLFile(ctx context.Context, serviceName string) (string, error) {
	base, err := baseDir()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s/config/%s", base, serviceName, postgresConfigFileName)
	slog.Debug(ctx, "Loading postgres sql file", map[string]string{
		"postgres_config_path": path,
	})

	c, err := ioutil.ReadFile(path)
	if err != nil {
		return "", terrors.Augment(err, "Failed to load config postgres sql file.", nil)
	}

	slog.Info(ctx, "Loaded config postgres sql file.")
	return string(c), nil
}

// baseDir returns the absolute path of the project root; which is the grandparent of this package.
func baseDir() (string, error) {
	c, err := os.Getwd()
	if err != nil {
		return "", terrors.Augment(err, "Failed to get current working directory", nil)
	}

	d := filepath.Dir(c)
	dd := filepath.Dir(d)

	baseDir, err := filepath.Abs(dd)
	if err != nil {
		return "", terrors.Augment(err, "Failed to get absolute parent dir of current working directory", nil)
	}

	return baseDir, nil
}
