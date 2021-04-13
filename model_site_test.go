package ambient_test

import (
	"testing"

	"github.com/josephspurrier/ambient"
	"github.com/stretchr/testify/assert"
)

func TestSiteURL(t *testing.T) {
	s := new(ambient.Site)
	s.Scheme = "http"
	s.URL = "localhost"
	assert.Equal(t, "http://localhost", s.SiteURL())
}
