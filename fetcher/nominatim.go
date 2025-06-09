package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
	"github.com/rs/zerolog/log"
)

// NominatimFetcher implements LocationFetcher using the Nominatim API.
type nominatimFetcher struct{}

// Creates a fetcher that geocodes locations using the Nominatim API.
func NewNomnatimFetcher() LocationFetcher {
	return &nominatimFetcher{}
}

// Fetch fetches locations from the Nominatim API for the given query.
func (f *nominatimFetcher) Fetch(query string) ([]location.Location, error) {

	log.Debug().Str("Nominatim query", query).Msg("Fetching location from Nominatim")

	req, err := buildNominatimRequest(query)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data []location.Location
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("cannot parse Nominatim response: %w\nThe response body was %s", err, string(body))
	}
	return data, nil
}

// buildNominatimRequest creates an HTTP GET request for the Nominatim API for the given query.
func buildNominatimRequest(query string) (*http.Request, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", url.QueryEscape(query))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "owen-feehan geocoding (owen@owenfeehan.com)")
	return req, nil
}

// Assert implementation
var _ LocationFetcher = (*nominatimFetcher)(nil)
