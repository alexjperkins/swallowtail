package multiplexing

import (
	"context"
	"sync"
)

// Multiplexer interface; defines the behaviour of an interface.
type Multiplexer interface {
	// Adds a consumer to the multiplexer
	AddConsumer(chan interface{}) (int, error)
	// Removes a consumer from the multiplexer
	RemoveConsumerByID(id int) error
	// Starts multiplexing
	Start(context.Context, chan interface{})
	// Stops multiplexing
	Stop()
}

// Multiplex multiplexs an input chan onto n consumers
type Multiplex struct {
	consumerGroup []*MultiplexConsumer
	consumerMtx   sync.RWMutex

	mtx  sync.Mutex
	done chan struct{}
}

func New(consumerGroup []*MultiplexConsumer) *Multiplex {
	return &Multiplex{
		consumerGroup: consumerGroup,
		done:          make(chan struct{}, 1),
	}

}

// Start multiplexing input into consumers
func (m *Multiplex) Start(ctx context.Context, input chan interface{}) {
	m.mtx.Lock()
	go func() {
		defer m.mtx.Unlock()
		for {
			select {
			case e := <-input:
				for _, c := range m.getConsumers() {
					c.send(ctx, e)

				}
			case <-ctx.Done():
				for _, c := range m.getConsumers() {
					c.close()
				}
				return
			case <-m.done:
				for _, c := range m.getConsumers() {
					c.close()
				}
				return
			}
		}
	}()
}

func (m *Multiplex) AddConsumer(consumer *MultiplexConsumer) (int, error) {
	m.consumerMtx.Lock()
	defer m.consumerMtx.Unlock()
	id := len(m.consumerGroup)
	m.consumerGroup = append(m.consumerGroup, consumer)
	return id, nil
}

// RemoveConsumer does nothing right now
func (m *Multiplex) RemoveConsumer(id int) error {
	// we need a hash map here
	// does nothing for now
	return nil
}

func (m *Multiplex) Stop() {
	m.done <- struct{}{}
	close(m.done)
}

func (m *Multiplex) getConsumers() []*MultiplexConsumer {
	m.consumerMtx.RLock()
	defer m.consumerMtx.RUnlock()
	return m.consumerGroup
}
