package datastorage

import (
	"encoding/json"

	"github.com/josephspurrier/ambient/app/lib/envdetect"
	"github.com/josephspurrier/ambient/app/model"
)

// Datastorer reads and writes data to an object.
type Datastorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

// Storage represents a writable and readable object.
type Storage struct {
	Site         *model.Site
	PluginRoutes *model.PluginRoutes
	datastorer   Datastorer
}

// New returns a writable and readable site object. Returns an error if the
// object cannot be initially read.
func New(ds Datastorer, site *model.Site) (*Storage, error) {
	s := &Storage{
		Site:       site,
		datastorer: ds,
	}

	err := s.Load()
	if err != nil {
		return nil, err
	}

	// Initialize the data structure.
	s.PluginRoutes = &model.PluginRoutes{
		Routes: make(map[string][]model.Route),
	}

	// Fill in the missing defaults.
	s.Site.Correct()

	return s, nil
}

// Save writes the site object to the data storage and returns an error if it
// cannot be written.
func (s *Storage) Save() error {
	var b []byte
	var err error

	if envdetect.RunningLocalDev() {
		// Indent so the data is easy to read.
		b, err = json.MarshalIndent(s.Site, "", "    ")
	} else {
		b, err = json.Marshal(s.Site)
	}

	if err != nil {
		return err
	}

	err = s.datastorer.Save(b)
	if err != nil {
		return err
	}

	return nil
}

// Load reads the site object from the data storage and returns an error if
// it cannot be read.
func (s *Storage) Load() error {
	b, err := s.datastorer.Load()
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, s.Site)
	if err != nil {
		return err
	}

	return nil
}
