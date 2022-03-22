package config

import (
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/amberror"
)

// SetTitle sets the title.
func (p *PluginSystem) SetTitle(title string) error {
	p.storage.site.Title = title
	return p.storage.Save()
}

// Title returns the title.
func (p *PluginSystem) Title() string {
	return p.storage.site.Title
}

// SetScheme sets the site scheme.
func (p *PluginSystem) SetScheme(scheme string) error {
	p.storage.site.Scheme = scheme
	return p.storage.Save()
}

// Scheme returns the site scheme.
func (p *PluginSystem) Scheme() string {
	return p.storage.site.Scheme
}

// SetURL sets the site URL.
func (p *PluginSystem) SetURL(URL string) error {
	p.storage.site.URL = URL
	return p.storage.Save()
}

// URL returns the URL without the scheme at the beginning.
func (p *PluginSystem) URL() string {
	return p.storage.site.URL
}

// FullURL returns the URL with the scheme at the beginning.
func (p *PluginSystem) FullURL() string {
	return p.storage.site.SiteURL()
}

// Updated returns the home last updated timestamp.
func (p *PluginSystem) Updated() time.Time {
	return p.storage.site.Updated
}

// Tags returns the list of tags.
func (p *PluginSystem) Tags(onlyPublished bool) ambient.TagList {
	return p.storage.site.Tags(onlyPublished)
}

// SetContent sets the home page content.
func (p *PluginSystem) SetContent(content string) error {
	p.storage.site.Content = content
	return p.storage.Save()
}

// Content returns the site home page content.
func (p *PluginSystem) Content() string {
	return p.storage.site.Content
}

// SavePost saves a post.
func (p *PluginSystem) SavePost(ID string, post ambient.Post) error {
	p.storage.site.Posts[ID] = post
	return p.storage.Save()
}

// PostsAndPages returns the list of posts and pages.
func (p *PluginSystem) PostsAndPages(onlyPublished bool) ambient.PostWithIDList {
	return p.storage.site.PostsAndPages(onlyPublished)
}

// PublishedPosts returns the list of published posts.
func (p *PluginSystem) PublishedPosts() []ambient.Post {
	return p.storage.site.PublishedPosts()
}

// PublishedPages returns the list of published pages.
func (p *PluginSystem) PublishedPages() []ambient.Post {
	return p.storage.site.PublishedPages()
}

// PostBySlug returns the post by slug.
func (p *PluginSystem) PostBySlug(slug string) ambient.PostWithID {
	return p.storage.site.PostBySlug(slug)
}

// PostByID returns the post by ID.
func (p *PluginSystem) PostByID(ID string) (ambient.Post, error) {
	post, ok := p.storage.site.Posts[ID]
	if !ok {
		return ambient.Post{}, amberror.ErrNotFound
	}

	return post, nil
}

// DeletePostByID deletes a post.
func (p *PluginSystem) DeletePostByID(ID string) error {
	delete(p.storage.site.Posts, ID)
	return p.storage.Save()
}
