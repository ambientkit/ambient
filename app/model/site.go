// Package model provides the data structure for the application.
package model

import (
	"fmt"
	"sort"
	"time"
)

// Site -
type Site struct {
	Title          string                    `json:"title"`
	Subtitle       string                    `json:"subtitle"`
	Footer         string                    `json:"footer"`
	Scheme         string                    `json:"scheme"`
	URL            string                    `json:"url"`
	LoginURL       string                    `json:"loginurl"`
	Created        time.Time                 `json:"created"`
	Updated        time.Time                 `json:"updated"`
	Content        string                    `json:"content"` // Home content.
	PluginSettings map[string]PluginSettings `json:"plugins"`
	PluginFields   map[string]PluginFields   `json:"pluginfields"`
	Posts          map[string]Post           `json:"posts"`
}

// Correct will fill in the missing defaults.
func (s *Site) Correct() {
	// Set the defaults for the site object.
	// Save to storage. Ensure the posts exists first so it doesn't error.
	if s.Posts == nil {
		s.Posts = make(map[string]Post)
	}
	if s.PluginSettings == nil {
		s.PluginSettings = make(map[string]PluginSettings)
	}
	if s.PluginFields == nil {
		s.PluginFields = make(map[string]PluginFields)
	}
	// Ensure redirects don't try to happen if the scheme is empty.
	if s.Scheme == "" {
		s.Scheme = "http"
	}
	// Ensure it's set to the login page works.
	if s.LoginURL == "" {
		s.LoginURL = "admin"
	}
}

// SiteURL -
func (s Site) SiteURL() string {
	return fmt.Sprintf("%v://%v", s.Scheme, s.URL)
}

// SiteTitle -
func (s Site) SiteTitle() string {
	return fmt.Sprintf("%v", s.Title)
}

// SiteSubtitle -
func (s Site) SiteSubtitle() string {
	return fmt.Sprintf("%v", s.Subtitle)
}

// PublishedPosts -
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

// PublishedPages -
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

// PostsAndPages -
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

// Tags -
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

// PostBySlug -
func (s Site) PostBySlug(slug string) PostWithID {
	// FIXME: This needs to be optimized.
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
