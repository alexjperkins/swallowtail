package cassandra

import (
	"fmt"
	"swallowtail/libraries/environment"

	"github.com/monzo/gocassa"
)

// KeyspaceWithName returns a new cassandra keyspace using the provided name.
func KeyspaceWithName(keyspace string, cfg environment.Cassandra) (gocassa.KeySpace, error) {
	if err := validateKeyspace(keyspace); err != nil {
		return nil, fmt.Errorf("create new keyspace with name: %w", err)
	}

	// Establish a connection.
	conn, err := gocassa.Connect(cfg.SeedNodeIPs, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("establish connection: %w", err)
	}

	return conn.KeySpace(keyspace), nil
}
