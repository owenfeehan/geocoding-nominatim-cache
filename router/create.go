package router

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// createRouter creates a GIN router with logging configured to write to zerolog.
//
// proxyList is a comma-separated list of trusted proxy IPs or CIDRs.
//
// It will exit with a fatal error if setting the proxy list fails.
func createRouter(proxyList string) (*gin.Engine, error) {
	// Write prettified console output with timestamps to stdout
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Create a root logger that writes to our consoleWriter
	rootLogger := zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	// Override Gin’s default writers
	gin.DefaultWriter = rootLogger
	gin.DefaultErrorWriter = rootLogger

	router := gin.New()
	// This will now emit to zerolog
	router.Use(gin.Logger()) // writes INFO-level request logs
	router.Use(gin.Recovery())

	if err := configureTrustedProxies(router, proxyList); err != nil {
		return nil, fmt.Errorf("failed to configure trusted proxies: %s for %w", proxyList, err)
	}

	return router, nil
}

// configureTrustedProxies configures the Gin router to trust specific proxies (or none if the list is empty).
func configureTrustedProxies(router *gin.Engine, proxyList string) error {
	// Convert the comma-separated string of trusted proxies into a slice of strings
	proxies := convertToArray(proxyList)

	// Set the trusted proxies for the Gin router
	return router.SetTrustedProxies(proxies)
}

// convertTrustedProxies converts a comma-separated string of trusted proxies into a slice of strings.
//
// nil is returned if proxyList is empty
func convertToArray(proxyList string) []string {
	if proxyList == "" {
		return nil
	}
	splitList := strings.Split(proxyList, ",")

	// Trim any whitespace from each element
	for i := range splitList {
		splitList[i] = strings.TrimSpace(splitList[i])
	}

	return splitList
}
