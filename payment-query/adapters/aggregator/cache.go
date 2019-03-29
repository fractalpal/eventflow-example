package aggregator

import (
	"sync"

	"github.com/fractalpal/eventflow-example/payment-query/app"
)

type PaymentsCache interface {
	Set(id string, payment app.Payment)
	Get(id string) *app.Payment
	Remove(id string)
}

// InMemoryCache for storing payments
type InMemoryCache struct {
	// for id aggregator
	data    map[string]app.Payment
	rwMutex sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: map[string]app.Payment{},
	}
}

func (c *InMemoryCache) Set(id string, payment app.Payment) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.data[id] = payment
}

func (c *InMemoryCache) Get(id string) *app.Payment {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	if p, ok := c.data[id]; ok {
		return &p
	}
	return nil
}

func (c *InMemoryCache) Remove(id string) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	delete(c.data, id)
}
