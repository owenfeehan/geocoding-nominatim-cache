package fetcher

import (
	"sync"
	"time"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
)

// Throttler wraps a LocationFetcher and ensures at most one request per second.
type throttler struct {
	delegate LocationFetcher

	// A minimum delay between calls to the delegate. The thread will sleep if necessary to ensure this delay.
	minDelay time.Duration

	mu       sync.Mutex
	lastCall time.Time
}

// NewThrottler creates a new Throttler that wraps the given delegate.
//
// the minDelay parameter is the minimum time to wait between calls to the delegate.
func NewThrottler(delegate LocationFetcher, minDelay time.Duration) LocationFetcher {
	return &throttler{delegate: delegate, minDelay: minDelay}
}

// Fetch calls the delegate's Fetch method, ensuring at most one call per second (thread-safe).
func (t *throttler) Fetch(query string) ([]location.Location, error) {
	t.mu.Lock()
	now := time.Now()
	wait := time.Second - now.Sub(t.lastCall)
	if wait > 0 {
		t.mu.Unlock()
		time.Sleep(wait)
		t.mu.Lock()
	}
	t.lastCall = time.Now()
	t.mu.Unlock()
	return t.delegate.Fetch(query)
}

// Assert implementation
var _ LocationFetcher = (*throttler)(nil)
