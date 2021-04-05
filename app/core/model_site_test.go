package core_test

import (
	"testing"

	"github.com/josephspurrier/ambient/app/core"
	"github.com/stretchr/testify/assert"
)

func TestSiteURL(t *testing.T) {
	s := new(core.Site)
	s.Scheme = "http"
	s.URL = "localhost"
	assert.Equal(t, "http://localhost", s.SiteURL())
}
