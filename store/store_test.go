package store

import (
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

func testWithStore(t *testing.T, store LocationStore) {
	key := store.BuildKey("Brussels")
	locs := []location.Location{{DisplayName: "Brussels, Belgium", Latitude: "50.8503", Longitude: "4.3517"}}

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
	if len(got) == 0 || got[0].DisplayName != locs[0].DisplayName {
		t.Errorf("Expected %v, got %v", locs[0], got)
	}
}
