package away_test

import (
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type route struct {
	pattern string
	method  string
	segs    []string
	handler http.Handler
	prefix  bool
}

type routeList []route

func (s routeList) Len() int {
	return len(s)
}

func (s routeList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s routeList) Less(i, j int) bool {
	var si string = s[i].pattern
	var sj string = s[j].pattern
	var siLower = strings.ToLower(si)
	var sjLower = strings.ToLower(sj)
	if strings.HasPrefix(sjLower, "/:") {
		return true
	} else if strings.HasPrefix(siLower, "/:") {
		return false
	} else if strings.Contains(sjLower, ":") && !strings.Contains(siLower, ":") {
		return true
	} else if !strings.Contains(sjLower, ":") && strings.Contains(siLower, ":") {
		return false
	}

	if siLower == sjLower {
		return si < sj
	}
	return siLower < sjLower
}

func TestMain(t *testing.T) {
	var arr routeList = []route{
		{
			method:  "GET",
			pattern: "/:slug",
		},
		{
			method:  "GET",
			pattern: "/",
		},
		{
			method:  "GET",
			pattern: "/cool/balloon",
		},
		{
			method:  "GET",
			pattern: "/cool/:slug",
		},
		{
			method:  "GET",
			pattern: "/cool/another",
		},
		{
			method:  "GET",
			pattern: "/rss/:ok",
		},
		{
			method:  "GET",
			pattern: "/cool",
		},
		{
			method:  "GET",
			pattern: "/rss.xml",
		},
	}

	assert.Equal(t, arr[0].pattern, "/:slug")
	sort.Sort(arr)
	assert.Equal(t, arr[0].pattern, "/")
	assert.Equal(t, arr[len(arr)-1].pattern, "/:slug")
}
