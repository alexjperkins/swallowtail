package streams

import "context"

// Event ...
type Event struct {
	context context.Context

	attempts   int
	paritition int32
	offset     int64

	Metadata     map[string]string `json:"metadata"`
	Payload      interface{}       `json:"payload"`
	PartitionKey string            `json:"partition_key"`
}
