package dao

import (
	"fmt"

	"swallowtail/libraries/cassandra"
	"swallowtail/libraries/environment"

	"github.com/monzo/gocassa"
)

var (
	keyspaceName = "bookmarker"
    Keyspace gocassa.KeySpace
)

// Init initializes the DAO.
func Init(cfg environment.Cassandra) error {
	if Keyspace != nil {
		return fmt.Errorf("cassandra keyspace already initialised")
	}

	// Get keyspace by name.
	var err error
	Keyspace, err = cassandra.KeyspaceWithName(keyspaceName, cfg)
	if err != nil {
		return fmt.Errorf("create cassandra keyspace: %w", err)
	}

	return nil
}
