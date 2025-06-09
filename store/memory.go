package store

import (
	"sync"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
	"github.com/rs/zerolog/log"
)

// memoryStore is an in-memory implementation of LocationStore using a map and mutex for thread safety.
type memoryStore struct {
	mu    sync.RWMutex                   // protects store
	store map[string][]location.Location // cache storage
}

// NewMemoryStore creates a new in-memory-only implementation of LocationStore
//
// It uses mutex for thread safety and a map to store the locations.
//
// No eviction occurs, and once a value is set it is guaranteed to remain.
func NewMemoryStore() LocationStore {
	log.Info().Msg("Using in-memory store for locations")
	return &memoryStore{
		store: make(map[string][]location.Location),
	}
}

// BuildKey returns the cache key for a given query. For memoryStore, this is just the query string.
func (c *memoryStore) BuildKey(query string) string {
	return query
}

// Get retrieves the cached locations for the given key, or nil if not found.
func (c *memoryStore) Get(key string) ([]location.Location, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.store[key]; ok {
		return val, nil
	}
	return nil, nil
}

// Set stores the locations in the cache under the given key.
func (c *memoryStore) Set(key string, locations []location.Location) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = locations
	return nil
}

// Close is a no-op for memoryStore.
func (c *memoryStore) Close() error {
	return nil
}

// Assert implementation of LocationStore interface.
var _ LocationStore = (*memoryStore)(nil)
