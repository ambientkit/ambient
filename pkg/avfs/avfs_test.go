package avfs_test

import (
	"testing"

	"github.com/ambientkit/ambient/pkg/avfs"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	fileContents := []byte("file contents")
	efs := avfs.NewFS()
	efs.AddFile("template/content/hello.tmpl", fileContents)

	fff, err := efs.Open("template/content/hello.tmpl")
	assert.NoError(t, err)

	fsi, err := fff.Stat()
	assert.NoError(t, err)

	assert.Equal(t, string(fileContents), string(fsi.Sys().([]byte)))
}
