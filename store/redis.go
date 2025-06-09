package store

import (
	"context"
	"fmt"

	"strings"

	"github.com/owenfeehan/geocoding-nominatim-cache/location"
	redis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

// redisStore implements LocationStore using Redis as the backend.
type redisStore struct {
	redis *redis.Client
}

// NewRedisStore creates a new LocationStore using Redis as a backend.
//
// The address should be in the form "host:port" (e.g., "localhost:6379").
//
// Persistence is not guaranteed, and evication may occus based on the max-memory settings in Redis e.g.
//
//	CONFIG SET maxmemory 100mb
//	CONFIG SET maxmemory-policy allkeys-lru
//
// These can be set in a presistent way in the redis.conf, see https://redis.io/docs/latest/operate/rs/databases/memory-performance/eviction-policy/
func NewRedisStore(address string) LocationStore {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	log.Info().Str("Redis store address", address).Msg("Connecting to Redis store")

	return &redisStore{redis: client}
}

func (c *redisStore) BuildKey(query string) string {
	return fmt.Sprintf("geocode:%s", strings.ToLower(query))
}

func (c *redisStore) Get(key string) ([]location.Location, error) {
	cached, err := c.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	result, err := unmarshalLocations([]byte(cached))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func (c *redisStore) Set(key string, locations []location.Location) error {
	body, err := marshalLocations(locations)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, body, 0).Err()
}

func (c *redisStore) Close() error {
	return c.redis.Close()
}

// Assert implementation
var _ LocationStore = (*redisStore)(nil)
