package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/owenfeehan/geocoding-nominatim-cache/fetcher"
	"github.com/owenfeehan/geocoding-nominatim-cache/router"
	"github.com/owenfeehan/geocoding-nominatim-cache/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title			Owen's Geocoding API
// @version		1.0
// @description	An API for caching geocoding locations as fetched from Nominatim.
//
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host		localhost:8080
// @BasePath	/
func main() {

	// START: Flags for command-line arguments

	// Router related-flags
	addr := flag.String("address", "localhost:8080", "The address to bind the server to")
	var proxyList string
	flag.StringVar(&proxyList, "trusted-proxies", "", "Comma-separated list of trusted proxy IPs or CIDRs")

	// Location store related-flags
	redis := flag.String("redis", "", "Binds to a redis server at the given address (e.g., localhost:6379). If not set, uses BadgerDB as the default store.")
	inMemory := flag.Bool("inMemory", false, "Uses in-memory (non-persistent) location storage, ignoring Redis or BadgerDB. This takes precedence over the redis flag, if set.")

	// Location fetcher related-flags
	throttle := flag.Int("throttle", 2000, "The minimum number of milli-seconds between requests to the Nominatim API. Default is 2000 milliseconds. This must be at least 1000 milliseconds to comply with the Nominatim API usage policy.")

	// other flags
	debug := flag.Bool("debug", false, "Enable debug logging and debug mode on the web server.")
	// ENDT: Flags for command-line arguments

	flag.Parse()

	configureLogging(*debug)

	// Create a store for locations, using Redis or BadgerDB, or in-memory storage.
	locStore, err := createStore(*redis, *inMemory)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create location-store")
		return
	}

	// Ensure the store is closed when the application exits
	defer locStore.Close()

	// Create a fetcher for locations, using Nominatim API with throttling.
	locFetcher, err := createFetcher(*throttle)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create location-fetcher")
		return
	}

	// Create a Gin router and configure it with the application routes
	appRoutes := app{
		Store:   locStore,
		Fetcher: locFetcher,
	}

	err = router.CreateRunRouter(*addr, proxyList, appRoutes.ForwardGeocode)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
	} else {
		log.Info().Str("address", *addr).Msg("Server is running")
	}
}

// configures the logging and debug-mode settings for the application.
func configureLogging(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode) // Set Gin to release mode for production
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
}

// Creates a store, using Redis (if redisAddr is non-empty) or otherwise BadgerDB.
func createStore(redisAddr string, inMemory bool) (store.LocationStore, error) {

	if inMemory {
		return store.NewMemoryStore(), nil
	}

	if redisAddr != "" {
		return store.NewRedisStore(redisAddr), nil
	} else {
		return store.NewBadgerStore(nil) // Use nil to automatically determine a path from the application-dir
	}
}

// Creates a fetcher for locations, using Nominatim API with throttling.
func createFetcher(throttle int) (fetcher.LocationFetcher, error) {

	if throttle < 1000 {
		return nil, fmt.Errorf("throttle must be at least 1000 milliseconds to comply with the Nominatim API usage policy")
	}

	log.Debug().Int("throttle duration in millis", throttle).Msg("Throttling Nomatim requests")

	// As per the Nominatim API usage policy, we should not send requests more frequently than once every 1 second.
	// We throttle to a 2 second delay to be conservative and avoid hitting the rate limit.
	// See https://operations.osmfoundation.org/policies/nominatim/
	throttle_duration := time.Duration(throttle) * time.Millisecond
	return fetcher.NewThrottler(fetcher.NewNomnatimFetcher(), throttle_duration), nil
}
