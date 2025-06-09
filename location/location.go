package location

// Location represents a Nominatim geocoding result.
//
// Latitude and longitude are stored as strings to maintain precision.
//
// The JSON names are deliberately chosen to match the Nominatim API response format.
type Location struct {
	DisplayName string `json:"display_name"`
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
}
