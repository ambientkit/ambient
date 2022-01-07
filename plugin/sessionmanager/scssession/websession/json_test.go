package websession_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/josephspurrier/ambient/lib/aesdata"
	"github.com/josephspurrier/ambient/plugin/sessionmanager/scssession/websession"
	"github.com/josephspurrier/ambient/plugin/storage/localstorage/store"
	"github.com/stretchr/testify/assert"
)

func TestNewJSONSession(t *testing.T) {
	// Use local filesytem when developing.
	f := "data.bin"
	err := ioutil.WriteFile(f, []byte(""), 0644)
	assert.NoError(t, err)
	ss := store.NewLocalStorage(f)

	// Set up the session storage provider.
	secretkey := "82a18fbbfed2694bb15d512a70c53b1a088e669966918d3d474564b2ac44349b"
	en := aesdata.NewEncryptedStorage(secretkey)
	store, err := websession.NewJSONSession(ss, en)
	assert.NoError(t, err)

	token := "abc"
	data := "hello"
	now := time.Now()

	err = store.Commit(token, []byte(data), now)
	assert.NoError(t, err)

	b, exists, err := store.Find(token)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, data, string(b))

	err = store.Delete(token)
	assert.NoError(t, err)

	_, exists, err = store.Find(token)
	assert.NoError(t, err)
	assert.False(t, exists)

	os.Remove(f)
}
