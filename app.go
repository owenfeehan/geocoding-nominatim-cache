package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/owenfeehan/geocoding-nominatim-cache/fetcher"
	"github.com/owenfeehan/geocoding-nominatim-cache/store"
)

// The global state needed across the handlers
// app holds the dependencies for the HTTP handlers.
type app struct {
	Store   store.LocationStore
	Fetcher fetcher.LocationFetcher
}

// Needed to document the error response for Swagger
type ErrorResponse struct {
	Error string `json:"error" example:"invalid input"`
}

// forwardGeocode handles the /locations/:place endpoint.
//
// @Summary      Get location coordinates for a placename
// @Description  get location coordinates and a canonical placename from a placename-query-string
// @Accept       json
// @Produce      json
// @Param        place   path      string  true  "query indicating a place or address"
// @Success      200  {object}  location.Location
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /locations/{place} [get]
func (a *app) ForwardGeocode(c *gin.Context) {
	place := c.Param("place")
	if place == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "place parameter is required"})
		return
	}

	loc, err := queryLocation(a.Store, a.Fetcher, place)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, loc)
}
