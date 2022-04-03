package ambient

// DataStorer reads and writes data to an object.
type DataStorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}

// // DataStorer reads and writes data to a store.
// type DataStorer interface {
// 	// Commit should add the data to the store. If the key already exists, then
// 	// the data and should be overwritten.
// 	Commit(collection string, key string, b []byte) (err error)

// 	// Find should return the data for a key from the store. If the key is not
// 	// found, the found return value should be false (and the err return value
// 	// should be nil). Similarly, malformed data should result in a found return
// 	// value of false and a nil err value. The err return value should be used
// 	// for system errors only.
// 	Find(collection string, key string) (b []byte, found bool, err error)

// 	// Delete should remove the key and corresponding data from the store. If
// 	// the key does not exist then Delete should be a no-op and return nil
// 	// (not an error).
// 	Delete(collection string, key string) (err error)
// }
