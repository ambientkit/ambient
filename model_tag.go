package ambient

import (
	"strings"
	"time"
)

// TagList represents a list of sortable tags.
type TagList []Tag

func (t TagList) Len() int {
	return len(t)
}
func (t TagList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t TagList) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// Tag represents a tag on a post or page.
type Tag struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

// String returns a comma separated list of tags.
func (t TagList) String() string {
	arr := make([]string, 0)
	for _, v := range t {
		arr = append(arr, v.Name)
	}

	return strings.Join(arr, ",")
}

// Split returns a list of tags from a comma separated list.
func (t TagList) Split(s string) TagList {
	trimmed := strings.TrimSpace(s)

	// Return an empty object since split returns 1 element when empty.
	if len(trimmed) == 0 {
		return TagList{}
	}

	ts := time.Now()

	arrTags := make([]Tag, 0)
	tags := strings.Split(trimmed, ",")
	for _, v := range tags {
		arrTags = append(arrTags, Tag{
			Name:      strings.TrimSpace(v),
			Timestamp: ts,
		})
	}

	return arrTags
}
