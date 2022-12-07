package environment

import (
	"errors"
)

const environmentFileName = "ENVIRONMENT_FILE"

var (
	ErrMissingEnvironmentFileEnvVar = errors.New("missing environment file environment variable")
)

// Environment defines the full environment for the application as a typed struct.
type Environment struct {
	Cassandra Cassandra
	Metadata  Metadata
}

type Metadata struct {
	ServiceName string `envconfig:"SERVICE_NAME"`
}

// Cassandra defines the cassandra specific environment.
type Cassandra struct {
	SeedNodeIPs []string `envconfig:"CASSANDRA_CONNECTION_URL"`
	Keyspace    string   `envconfig:"CASSANDRA_KEYSPACE" default:"birdperch"`
	Username    string   `envconfig:"CASSANDRA_USERNAME"`
	Password    string   `envconfig:"CASSANDRA_PASSWORD"`
}
