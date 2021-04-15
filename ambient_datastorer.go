package ambient

// IDataStorer reads and writes data to an object.
type IDataStorer interface {
	Save([]byte) error
	Load() ([]byte, error)
}
