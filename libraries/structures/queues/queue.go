package queues

// Queue interface
type Queue interface {
	Pop() (interface{}, bool)
	Peek() (interface{}, bool)
	Push(interface{}) error
	Len() int
	AtCapacity() bool
	// GetAsArray returns a copy of internal array of the queue; does not allow for mutation.
	GetAsArray() []interface{}
}
