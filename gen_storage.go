// Code generated by ifacemaker. DO NOT EDIT.

package ambient

// Storage provides app config functions.
type Storage interface {
	// Save writes the site object to the data storage and returns an error if it
	// cannot be written.
	Save() error
	// SaveDecrypted writes the site object to the data storage always decrypted and
	// returns an error if it cannot be written.
	SaveDecrypted() error
	// Load reads the site object from the data storage and returns an error if
	// it cannot be read.
	Load() error
	// LoadDecrypted reads the site object from the data storage always decrypted
	// and returns an error if it cannot be read.
	LoadDecrypted() error
}
