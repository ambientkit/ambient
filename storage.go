package ambient

import (
	"encoding/json"
	"time"

	"github.com/josephspurrier/ambient/lib/envdetect"
)

// Storage represents a writable and readable object.
type Storage struct {
	site       *Site
	datastorer DataStorer
}

// NewStorage returns a writable and readable site object. Returns an error if the
// object cannot be initially read.
func NewStorage(ds DataStorer) (*Storage, error) {
	s := &Storage{
		site:       &Site{},
		datastorer: ds,
	}

	err := s.Load()
	if err != nil {
		return nil, err
	}

	// Fill in the missing defaults.
	s.site.Correct()

	return s, nil
}

// Save writes the site object to the data storage and returns an error if it
// cannot be written.
func (s *Storage) Save() error {
	var b []byte
	var err error

	// Save the updated timestamp.
	s.site.Updated = time.Now()

	if envdetect.RunningLocalDev() {
		// Indent so the data is easy to read.
		b, err = json.MarshalIndent(s.site, "", "    ")
	} else {
		b, err = json.Marshal(s.site)
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

	err = json.Unmarshal(b, s.site)
	if err != nil {
		return err
	}

	return nil
}
