package fetcher

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
)

type mockFetcher struct {
	calls int32
	delay time.Duration
}

func (m *mockFetcher) Fetch(query string) ([]location.Location, error) {
	atomic.AddInt32(&m.calls, 1)
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return []location.Location{{DisplayName: query}}, nil
}

func TestThrottler_RespectsMinDelay(t *testing.T) {
	mock := &mockFetcher{}
	throttler := NewThrottler(mock, 200*time.Millisecond)

	start := time.Now()
	_, _ = throttler.Fetch("A")
	_, _ = throttler.Fetch("B")

	assertMinDuration(t, start)
	assertCalls(t, mock, 2)
}

func TestThrottler_Concurrent(t *testing.T) {
	mock := &mockFetcher{}
	throttler := NewThrottler(mock, 100*time.Millisecond)

	start := time.Now()
	ch := make(chan struct{})

	// Call fetch and send a message to the channel when done.
	for i := range 3 {
		go func(idx int) {
			query := fmt.Sprintf("A%d", idx)
			_, _ = throttler.Fetch(query)
			ch <- struct{}{}
		}(i)
	}

	// Wait for all goroutines to finish.
	for range 3 {
		<-ch
	}

	assertMinDuration(t, start)
	assertCalls(t, mock, 3)
}

// Asserts the expected number of calls occurred on the mock fetcher.
func assertCalls(t *testing.T, mock *mockFetcher, want int32) {
	if got := atomic.LoadInt32(&mock.calls); got != want {
		t.Errorf("expected %d calls to delegate, got %d", want, got)
	}
}

// Asserts that the duration since start is at least the expected minimum.
func assertMinDuration(t *testing.T, start time.Time) {
	min := 200 * time.Millisecond
	dur := time.Since(start)
	if dur < min {
		t.Errorf("expected at least %v for the calls, got %v", min, dur)
	}
}
