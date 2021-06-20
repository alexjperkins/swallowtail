package dao

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgconn"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	c, err := Init(ctx, &pgconn.Config{
		Database: "swallowtail_test",
		User:     "test",
		Password: "test",
		Host:     "postgres",
		// Test port
		Port: 5433,
	})
	if err != nil {
		log.Fatalf("Failed to established connection to test db: %v", err)
	}

	code := m.Run()

	// Close database connection & cleanup.
	db.Exec(ctx, `DROP TABLE IF EXISTS accounts`)
	c()
	os.Exit(code)
}
