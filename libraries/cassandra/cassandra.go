package cassandra

import (
	"fmt"
	"swallowtail/libraries/environment"

	"github.com/hailocab/gocassa"
)

// Cassandra ...
type Cassandra interface {
	Table()
}

type defaultCassandraImpl struct {
	keyspace gocassa.KeySpace
}

func (d *defaultCassandraImpl) Table() {}

type mockCassandraImpl struct {
	keyspace gocassa.KeySpace
}

func (m *mockCassandraImpl) Table() {}

func New(cfg environment.Cassandra) (Cassandra, error) {
	keyspace, err := gocassa.ConnectToKeySpace(cfg.Keyspace, cfg.SeedNodeIPs, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("create new cassandra client: %w", err)
	}

	return &defaultCassandraImpl{
		keyspace: keyspace,
	}, nil
}

// NewMock returns a mocked cassandra keyspace.
func NewMock() Cassandra {
	return &mockCassandraImpl{gocassa.NewMockKeySpace()}
}
