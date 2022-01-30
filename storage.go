package ambient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ambientkit/ambient/pkg/envdetect"
)

// Storage represents a writable and readable object.
type Storage struct {
	log        AppLogger
	site       *Site
	datastorer DataStorer
	secure     StorageEncryption
}

// NewStorage returns a writable and readable site object. Returns an error if the
// object cannot be initially read.
func NewStorage(log AppLogger, ds DataStorer, es StorageEncryption) (*Storage, error) {
	s := &Storage{
		log:        log,
		site:       &Site{},
		datastorer: ds,
		secure:     es,
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
	return s.save(true)
}

// SaveDecrypted writes the site object to the data storage always decrypted and
// returns an error if it cannot be written.
func (s *Storage) SaveDecrypted() error {
	return s.save(false)
}

// save writes the site object to the data storage and returns an error if it
// cannot be written.
func (s *Storage) save(forceEncryption bool) error {
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

	// Encrypt if set.
	if s.secure != nil && forceEncryption {
		b, err = s.secure.Encrypt(b)
		if err != nil {
			return fmt.Errorf("ambient: could not encrypt storage data: %v", err.Error())
		}
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
	return s.load(false)
}

// LoadDecrypted reads the site object from the data storage always decrypted
// and returns an error if it cannot be read.
func (s *Storage) LoadDecrypted() error {
	return s.load(true)
}

// Load reads the site object from the data storage and returns an error if
// it cannot be read.
func (s *Storage) load(allowDecrypted bool) error {
	b, err := s.datastorer.Load()
	if err != nil {
		return err
	}

	// Decrypt if set.
	if s.secure != nil {
		if string(b) == "{}" {
			s.log.Debug("ambient: found new storage data file, will encrypt on save")
		} else {
			decrypted, err := s.secure.Decrypt(b)
			if err != nil {
				// Could be a new file so don't fail, but it will force encryption.
				if string(b) == "{}" {
					s.log.Info("ambient: found new storage data file, will encrypt on save")
				} else {
					if !allowDecrypted {
						return fmt.Errorf("ambient: could not decrypt storage data: %v", err.Error())
					}
					decrypted = b
				}
			}
			b = decrypted
		}
	}

	err = json.Unmarshal(b, s.site)
	if err != nil {
		return err
	}

	return nil
}
