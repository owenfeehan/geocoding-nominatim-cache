package main

import (
	"errors"
	"testing"

	"github.com/owenfeehan/geocoding-nominatim-cache/fetcher"
	"github.com/owenfeehan/geocoding-nominatim-cache/location"
	"github.com/owenfeehan/geocoding-nominatim-cache/store"
)

const testQuery = "Brussels"

type mockStore struct {
	store.LocationStore
	getFunc func(string) ([]location.Location, error)
	setFunc func(string, []location.Location) error
}

func (m *mockStore) Get(key string) ([]location.Location, error) {
	return m.getFunc(key)
}
func (m *mockStore) Set(key string, locs []location.Location) error {
	if m.setFunc != nil {
		return m.setFunc(key, locs)
	}
	return nil
}
func (m *mockStore) BuildKey(query string) string { return query }
func (m *mockStore) Close() error                 { return nil }

type mockFetcher struct {
	fetchFunc func(string) ([]location.Location, error)
}

func (m *mockFetcher) Fetch(query string) ([]location.Location, error) {
	return m.fetchFunc(query)
}

func TestQueryLocationCacheHit(t *testing.T) {
	want := location.Location{DisplayName: testQuery}
	store := &mockStore{
		getFunc: func(_ string) ([]location.Location, error) {
			return []location.Location{want}, nil
		},
	}
	fetcher := &mockFetcher{
		fetchFunc: func(_ string) ([]location.Location, error) {
			t.Fatal("fetcher should not be called on cache hit")
			return nil, nil
		},
	}
	got, err := queryLocation(store, fetcher, testQuery)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DisplayName != want.DisplayName {
		t.Errorf("expected %v, got %v", want.DisplayName, got.DisplayName)
	}
}

func TestQueryLocationCacheMissAndFetch(t *testing.T) {
	want := location.Location{DisplayName: testQuery}
	store := &mockStore{
		getFunc: func(key string) ([]location.Location, error) {
			return nil, nil
		},
		setFunc: func(key string, locs []location.Location) error {
			if len(locs) == 0 || locs[0].DisplayName != want.DisplayName {
				t.Errorf("expected to cache %v, got %v", want.DisplayName, locs)
			}
			return nil
		},
	}
	fetcher := &mockFetcher{
		fetchFunc: func(query string) ([]location.Location, error) {
			return []location.Location{want}, nil
		},
	}
	got, err := queryLocation(store, fetcher, testQuery)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.DisplayName != want.DisplayName {
		t.Errorf("expected %v, got %v", want.DisplayName, got.DisplayName)
	}
}

func TestQueryLocationErrorCases(t *testing.T) {
	storeErr := errors.New("store error")
	fetchErr := errors.New("fetch error")

	// Cache error
	store := &mockStore{
		getFunc: func(key string) ([]location.Location, error) {
			return nil, storeErr
		},
	}
	fetcher := &mockFetcher{}
	_, err := queryLocation(store, fetcher, testQuery)
	if err == nil || err.Error() != "failed to retrieve from the cache: "+storeErr.Error() {
		t.Errorf("expected store error, got %v", err)
	}

	// Fetch error
	store = &mockStore{
		getFunc: func(key string) ([]location.Location, error) {
			return nil, nil
		},
	}
	fetcher = &mockFetcher{
		fetchFunc: func(query string) ([]location.Location, error) {
			return nil, fetchErr
		},
	}
	_, err = queryLocation(store, fetcher, testQuery)
	if err == nil || err.Error() != "failed to fetch location: "+fetchErr.Error() {
		t.Errorf("expected fetch error, got %v", err)
	}

	// No locations found
	store = &mockStore{
		getFunc: func(key string) ([]location.Location, error) {
			return nil, nil
		},
		setFunc: func(key string, locs []location.Location) error { return nil },
	}
	fetcher = &mockFetcher{
		fetchFunc: func(query string) ([]location.Location, error) {
			return nil, nil
		},
	}
	_, err = queryLocation(store, fetcher, testQuery)
	if err == nil || err.Error() != "no locations found for query: "+testQuery {
		t.Errorf("expected no locations found error, got %v", err)
	}
}

// Assert implementation of mocked interfaces.
var _ store.LocationStore = (*mockStore)(nil)
var _ fetcher.LocationFetcher = (*mockFetcher)(nil)
