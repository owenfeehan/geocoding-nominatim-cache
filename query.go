package main

import (
	"fmt"

	"github.com/owenfeehan/geocoding-nominatim-cache/fetcher"
	"github.com/owenfeehan/geocoding-nominatim-cache/location"
	"github.com/owenfeehan/geocoding-nominatim-cache/store"
)

// queryLocation retrieves a location for the given query, using cache if possible.
func queryLocation(locStore store.LocationStore, locFettcher fetcher.LocationFetcher, query string) (location.Location, error) {
	cacheKey := locStore.BuildKey(query)

	// Try to get location from cache
	loc, err := locStore.Get(cacheKey)
	if err != nil {
		return location.Location{}, fmt.Errorf("failed to retrieve from the cache: %w", err)
	} else if loc != nil {
		// Successful cache hit
		return extractFirstLocation(loc, query)
	}

	// If not cached, fetch from Nominatim API
	loc, err = locFettcher.Fetch(query)
	if err != nil {
		return location.Location{}, fmt.Errorf("failed to fetch location: %w", err)
	}

	// Cache the result
	if err := locStore.Set(cacheKey, loc); err != nil {
		fmt.Println("Cache Error, could not cache: ", err)
	}

	return extractFirstLocation(loc, query)
}

// extractFirstLocation returns the first location from the slice or an error if empty.
func extractFirstLocation(data []location.Location, query string) (location.Location, error) {
	if len(data) == 0 {
		return location.Location{}, fmt.Errorf("no locations found for query: %s", query)
	}
	return data[0], nil
}
