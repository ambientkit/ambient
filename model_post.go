package ambient

import (
	"time"
)

// Post -
type Post struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Canonical string    `json:"canonical"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	Page      bool      `json:"page"`
	Tags      TagList   `json:"tags"`
}

// PostWithID -
type PostWithID struct {
	Post
	ID string `json:"id"`
}

// PostWithIDList -
type PostWithIDList []PostWithID

func (t PostWithIDList) Len() int {
	return len(t)
}
func (t PostWithIDList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t PostWithIDList) Less(i, j int) bool {
	if t[i].Timestamp.Equal(t[j].Timestamp) {
		return t[i].Title > t[j].Title // Sort by title ASC
	} else if t[i].Timestamp.Before(t[j].Timestamp) {
		return true // Sort by timestamp, DESC
	}

	return false
}

// PostList -
type PostList []Post

func (t PostList) Len() int {
	return len(t)
}
func (t PostList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t PostList) Less(i, j int) bool {
	if t[i].Timestamp.Equal(t[j].Timestamp) {
		return t[i].Title > t[j].Title // Sort by title ASC
	} else if t[i].Timestamp.Before(t[j].Timestamp) {
		return true // Sort by timestamp, DESC
	}

	return false
}
