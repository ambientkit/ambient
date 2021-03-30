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
	//return s[i].pattern < s[j].pattern
	var si string = s[i].pattern
	var sj string = s[j].pattern
	var siLower = strings.ToLower(si)
	var sjLower = strings.ToLower(sj)
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
			pattern: "/cool",
		},
		{
			method:  "GET",
			pattern: "/rss.xml",
		},
	}

	assert.Equal(t, arr[0].pattern, "/:slug")
	sort.Sort(sort.Reverse(arr))
	assert.Equal(t, arr[2].pattern, "/:slug")
}
