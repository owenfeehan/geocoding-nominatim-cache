package store

import (
	"reflect"
	"testing"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
)

func TestMemoryStore(t *testing.T) {
	testWithStore(t, NewMemoryStore())
}

func TestBadgerStore(t *testing.T) {
	path := t.TempDir()
	store, err := NewBadgerStore(&path)
	if err != nil {
		t.Fatalf("Failed to create Badger store: %v", err)
	}
	defer store.Close()
	testWithStore(t, store)
}

// testWithStore tests the provided LocationStore implementation by performing a series of queries and checking the results.
func testWithStore(t *testing.T, store LocationStore) {
	// Query that returns no locations
	testLocation(t, store, "Unknown place", []location.Location{})

	// Query that returns a single location
	locationBelgium := location.Location{DisplayName: "Brussels, Belgium", Latitude: "50.8503", Longitude: "4.3517"}
	testLocation(t, store, "Brussels", []location.Location{locationBelgium})

	// Query that returns two locations
	locationWisconsin := location.Location{DisplayName: "Brussels, Wisconsin", Latitude: "10.8503", Longitude: "14.3517"}
	testLocation(t, store, "Brussels", []location.Location{locationBelgium, locationWisconsin})
}

// testLocation tests the LocationStore implementation by storing and retrieving locations for a given query.
func testLocation(t *testing.T, store LocationStore, query string, locs []location.Location) {
	key := store.BuildKey(query)

	// Store locations in the cache
	err := store.Set(key, locs)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Retrieve locations from the cache
	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !reflect.DeepEqual(got, locs) {
		t.Errorf("Expected %v, got %v", locs[0], got)
	}
}
