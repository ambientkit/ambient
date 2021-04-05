package templatebuffer_test

import (
	"testing"

	"github.com/josephspurrier/ambient/app/lib/templatebuffer"
	"github.com/stretchr/testify/assert"
)

var body = `This is a test. {{.SiteURL}}`

func TestMain(t *testing.T) {
	data := map[string]interface{}{"SiteURL": "ok"}
	s, err := templatebuffer.ParseTemplate(body, nil, data)
	assert.NoError(t, err)

	assert.True(t, len(s) > 0)
}
