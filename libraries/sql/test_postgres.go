package sql

// SetPostgresConnectionURL sets the connection URL to be used; this should only really
// be used for testing purposes.
func SetPostgresConnectionURL(url string) {
	mu.Lock()
	defer mu.Unlock()

	connURL = "postgresql://test:test@postgres:5433/swallowtail_test"
	if url != "" {
		connURL = url
	}
}
