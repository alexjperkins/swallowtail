package dao

import (
	"context"
	"flag"
	"log"
	"os"
	"swallowtail/libraries/sql"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(0)
	}
	ctx := context.Background()

	connectionURL := os.Getenv("SWALLOWTAIL_TEST_POSTGRES_CONNECTION_URL")
	sql.SetPostgresConnectionURL(connectionURL)

	err := Init(ctx, "s.account")
	if err != nil {
		log.Fatalf("Failed to established connection to test db: %v", err)
	}

	// Run our test assertions.
	code := m.Run()

	// Close database connection & cleanup.
	db.Exec(ctx, `DROP TABLE IF EXISTS s_account_accounts CASCADE`)
	os.Exit(code)
}
