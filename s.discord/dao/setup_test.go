// +build intergration

package dao

import (
	"context"
	"log"
	"os"
	"swallowtail/libraries/sql"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	connectionURL := os.Getenv("SWALLOWTAIL_TEST_POSTGRES_CONNECTION_URL")
	sql.SetPostgresConnectionURL(connectionURL)

	err := Init(ctx, "s.discord")
	if err != nil {
		log.Fatalf("Failed to established connection to test db: %v", err)
	}

	// Cleanup before starting.
	db.Exec(ctx, `DROP TABLE IF EXISTS accounts`)

	// Run our test assertions.
	code := m.Run()

	// Close database connection & cleanup.
	db.Exec(ctx, `DROP TABLE IF EXISTS accounts`)
	os.Exit(code)
}
