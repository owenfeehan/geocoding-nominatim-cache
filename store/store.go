// Package store provides a key-value backend for storing keys (queries) and corresponding location-data.
//
// This data is typicallly formed elsewhere after geocoding a query (i.e. place name) into a list of locations.
//
// The package exposes an interface LocationStore that defines methods for getting and setting the location-data.
//
// Different implementations can be instantiated, as per the requirements of the application, such as in-memory caching or
// persistent storage.
package store

import "github.com/owenfeehan/geocoding-nominatim-cache/location"

// LocationStore defines the interface for getting and putting data in the cache-backend.
type LocationStore interface {
	// Translates a query into a cache-key (which is used for subsequent set/get operations)
	BuildKey(query string) string

	// Stores a location-values for a given key
	Set(key string, value []location.Location) error

	// Retrieves a location-values for a given key
	Get(key string) ([]location.Location, error)

	// Closes the store and releases any resources (no-op for in-memory)
	Close() error
}
