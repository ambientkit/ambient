package ambient

import "time"

// SetTitle sets the title.
func (p *PluginSystem) SetTitle(title string) error {
	p.storage.Site.Title = title
	return p.storage.Save()
}

// Title returns the title.
func (p *PluginSystem) Title() string {
	return p.storage.Site.Title
}

// SetScheme sets the site scheme.
func (p *PluginSystem) SetScheme(scheme string) error {
	p.storage.Site.Scheme = scheme
	return p.storage.Save()
}

// Scheme returns the site scheme.
func (p *PluginSystem) Scheme() string {
	return p.storage.Site.Scheme
}

// SetURL sets the site URL.
func (p *PluginSystem) SetURL(URL string) error {
	p.storage.Site.URL = URL
	return p.storage.Save()
}

// URL returns the URL without the scheme at the beginning.
func (p *PluginSystem) URL() string {
	return p.storage.Site.URL
}

// FullURL returns the URL with the scheme at the beginning.
func (p *PluginSystem) FullURL() string {
	return p.storage.Site.SiteURL()
}

// Updated returns the home last updated timestamp.
func (p *PluginSystem) Updated() time.Time {
	return p.storage.Site.Updated
}

// Tags returns the list of tags.
func (p *PluginSystem) Tags(onlyPublished bool) TagList {
	return p.storage.Site.Tags(onlyPublished)
}

// SetContent sets the home page content.
func (p *PluginSystem) SetContent(content string) error {
	p.storage.Site.Content = content
	return p.storage.Save()
}

// Content returns the site home page content.
func (p *PluginSystem) Content() string {
	return p.storage.Site.Content
}

// SavePost saves a post.
func (p *PluginSystem) SavePost(ID string, post Post) error {
	p.storage.Site.Posts[ID] = post
	return p.storage.Save()
}

// PostsAndPages returns the list of posts and pages.
func (p *PluginSystem) PostsAndPages(onlyPublished bool) PostWithIDList {
	return p.storage.Site.PostsAndPages(onlyPublished)
}

// PublishedPosts returns the list of published posts.
func (p *PluginSystem) PublishedPosts() []Post {
	return p.storage.Site.PublishedPosts()
}

// PublishedPages returns the list of published pages.
func (p *PluginSystem) PublishedPages() []Post {
	return p.storage.Site.PublishedPages()
}

// PostBySlug returns the post by slug.
func (p *PluginSystem) PostBySlug(slug string) PostWithID {
	return p.storage.Site.PostBySlug(slug)
}

// PostByID returns the post by ID.
func (p *PluginSystem) PostByID(ID string) (Post, error) {
	post, ok := p.storage.Site.Posts[ID]
	if !ok {
		return Post{}, ErrNotFound
	}

	return post, nil
}

// DeletePostByID deletes a post.
func (p *PluginSystem) DeletePostByID(ID string) error {
	delete(p.storage.Site.Posts, ID)
	return p.storage.Save()
}
