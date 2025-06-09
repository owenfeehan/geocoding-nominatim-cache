// Package fetcher provides an interface for fetching locations based on a query
//
// This may or may not involves calling an external geocoding API.
package fetcher

import location "github.com/owenfeehan/geocoding-nominatim-cache/location"

// LocationFetcher is a polymorphic interface for fetching locations from a query.
//
// Example:
//  NewNomnatimFetcher().Fetch("Galway, Ireland")
type LocationFetcher interface {
	Fetch(query string) ([]location.Location, error)
}
