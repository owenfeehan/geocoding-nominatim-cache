package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/dgraph-io/badger/v4"
	"github.com/getlantern/appdir"
	location "github.com/owenfeehan/geocoding-nominatim-cache/location"
)

// badgerStore implements LocationStore using BadgerDB as the backend.
type badgerStore struct {
	db *badger.DB
}

// NewBadgerStore opens (or creates) a BadgerDB at the given path and returns a LocationStore.
//
// This provides a convenient persistent data-store on the file-system.
//
// However, it does not perform cache expiration. It is recommended to use the Redis backend if this feature is required.
//
// If the path is nil, it will create a default data directory under the user's app data directory, otherwise it will use the
// provided folder-path for the badger DB.
func NewBadgerStore(path *string) (LocationStore, error) {

	// Calculate a path, if not already provided
	if path == nil {
		pathDataDir, err := pathDataDirectory()
		if err != nil {
			return nil, err
		}

		path = &pathDataDir
	}

	log.Info().Str("BadgerDB store directory", *path).Msg("Opening BadgerDB store")

	db, err := badger.Open(badger.DefaultOptions(*path).WithLoggingLevel(badger.WARNING))
	if err != nil {
		return nil, err
	}
	return &badgerStore{db: db}, nil
}

// Determines a path to where the BadgerDB data is stored (created if not already existing)
func pathDataDirectory() (string, error) {
	baseDir := appdir.General("geocoding")

	subDir := filepath.Join(baseDir, "locations")

	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		return subDir, fmt.Errorf("failed to create badger subdirectory: %w", err)
	}

	return subDir, nil
}

func (b *badgerStore) BuildKey(query string) string {
	return fmt.Sprintf("geocode:%s", query)
}

func (b *badgerStore) Get(key string) ([]location.Location, error) {
	var result []location.Location
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			res, err := unmarshalLocations(val)
			if err != nil {
				return err
			}
			result = res
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func (b *badgerStore) Set(key string, locations []location.Location) error {
	val, err := marshalLocations(locations)
	if err != nil {
		return err
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), val)
	})
}

// Close closes the underlying BadgerDB.
func (b *badgerStore) Close() error {
	return b.db.Close()
}

// Assert implementation
var _ LocationStore = (*badgerStore)(nil)
