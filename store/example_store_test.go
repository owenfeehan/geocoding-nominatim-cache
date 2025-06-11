package store

import (
	"fmt"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
)

func ExampleLocationStore() {
	locations := []location.Location{{DisplayName: "Brussels, Belgium", Latitude: "50.8503", Longitude: "4.3517"}}

	store := NewMemoryStore()
	key := store.BuildKey("Brussels") // Find a key corresponding to the query Use query string directly as key

	// Store locations
	if err := store.Set(key, locations); err != nil {
		fmt.Println("Set failed:", err)
		return
	}

	// Retrieve locations
	got, err := store.Get(key)
	if err != nil {
		fmt.Println("Get failed:", err)
		return
	}
	if len(got) > 0 {
		fmt.Printf("Found: %s (Lat: %s, Lon: %s)\n", got[0].DisplayName, got[0].Latitude, got[0].Longitude)
	} else {
		fmt.Println("No location found")
	}

	// Output:
	// Found: Brussels, Belgium (Lat: 50.8503, Lon: 4.3517)
}
