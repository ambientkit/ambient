package ambient

import (
	"fmt"
	"sort"
	"time"
)

// Site represents the site information that is in storage.
type Site struct {
	Title         string                `json:"title"`   // Title of the site.
	Content       string                `json:"content"` // Home or default content.
	Scheme        string                `json:"scheme"`  // http or https
	URL           string                `json:"url"`     // URL without scheme and without trailing slash.
	Updated       time.Time             `json:"updated"` // Save time the data was saved (not only changed).
	Posts         map[string]Post       `json:"posts"`   // List of posts.
	PluginStorage map[string]PluginData `json:"plugins"` // List of plugins, whether they are found, enabled, and what fields they support.
}

// PluginData represents the plugin storage information.
type PluginData struct {
	Enabled  bool           `json:"enabled"`
	Version  string         `json:"version"`
	Grants   PluginGrants   `json:"grants"`
	Settings PluginSettings `json:"settings"`
}

// Correct will fill in the missing defaults.
func (s *Site) Correct() {
	// Set the defaults for the site object.
	// Save to storage. Ensure the posts exists first so it doesn't error.
	if s.Posts == nil {
		s.Posts = make(map[string]Post)
	}
	if s.PluginStorage == nil {
		s.PluginStorage = make(map[string]PluginData)
	}
}

// SiteURL returns the URL with the scheme.
func (s Site) SiteURL() string {
	return fmt.Sprintf("%v://%v", s.Scheme, s.URL)
}

// PublishedPosts returns published posts (no pages).
func (s Site) PublishedPosts() []Post {
	arr := make(PostList, 0)
	for _, v := range s.Posts {
		if v.Published && !v.Page {
			arr = append(arr, v)
		}
	}

	sort.Sort(sort.Reverse(arr))

	return arr
}

// PublishedPages returns published pages (no posts).
func (s Site) PublishedPages() []Post {
	arr := make(PostList, 0)
	for _, v := range s.Posts {
		if v.Published && v.Page {
			arr = append(arr, v)
		}
	}

	sort.Sort(sort.Reverse(arr))

	return arr
}

// PostsAndPages returns list of posts and pages with IDs.
func (s Site) PostsAndPages(onlyPublished bool) PostWithIDList {
	arr := make(PostWithIDList, 0)
	for k, v := range s.Posts {
		if onlyPublished && !v.Published {
			continue
		}

		p := PostWithID{Post: v, ID: k}
		arr = append(arr, p)
	}

	sort.Sort(sort.Reverse(arr))

	return arr
}

// Tags returns a list of tags.
func (s Site) Tags(onlyPublished bool) TagList {
	// Get unique values.
	m := make(map[string]Tag)
	for _, v := range s.Posts {
		if onlyPublished && !v.Published {
			continue
		}

		for _, t := range v.Tags {
			m[t.Name] = t
		}
	}

	// Create unsorted tag list.
	arr := make(TagList, 0)
	for _, v := range m {
		arr = append(arr, v)
	}

	// Sort by name.
	sort.Sort(arr)

	return arr
}

// PostBySlug returns a post by slug/URL.
func (s Site) PostBySlug(slug string) PostWithID {
	// TODO: This needs to be optimized.
	var p PostWithID
	for k, v := range s.Posts {
		if v.URL == slug {
			p = PostWithID{
				Post: v,
				ID:   k,
			}
			break
		}
	}

	return p
}
