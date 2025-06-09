package store

import (
	"encoding/json"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
)

// marshalLocations serializes a slice of Location to JSON.
func marshalLocations(locs []location.Location) ([]byte, error) {
	return json.Marshal(locs)
}

// unmarshalLocations deserializes JSON data into a slice of Location.
func unmarshalLocations(data []byte) ([]location.Location, error) {
	var result []location.Location
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
