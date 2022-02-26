package secureconfig

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/config"
)

// SavePost saves a post.
func (ss *SecureSite) SavePost(ID string, post ambient.Post) error {
	if !ss.Authorized(ambient.GrantSitePostWrite) {
		return config.ErrAccessDenied
	}

	return ss.pluginsystem.SavePost(ID, post)
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(onlyPublished bool) (ambient.PostWithIDList, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, config.ErrAccessDenied
	}

	return ss.pluginsystem.PostsAndPages(onlyPublished), nil
}

// PublishedPosts returns the list of published posts.
func (ss *SecureSite) PublishedPosts() ([]ambient.Post, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, config.ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPosts(), nil
}

// PublishedPages returns the list of published pages.
func (ss *SecureSite) PublishedPages() ([]ambient.Post, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, config.ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPages(), nil
}

// PostBySlug returns the post by slug.
func (ss *SecureSite) PostBySlug(slug string) (ambient.PostWithID, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return ambient.PostWithID{}, config.ErrAccessDenied
	}

	return ss.pluginsystem.PostBySlug(slug), nil
}

// PostByID returns the post by ID.
func (ss *SecureSite) PostByID(ID string) (ambient.Post, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return ambient.Post{}, config.ErrAccessDenied
	}

	return ss.pluginsystem.PostByID(ID)
}

// DeletePostByID deletes a post.
func (ss *SecureSite) DeletePostByID(ID string) error {
	if !ss.Authorized(ambient.GrantSitePostDelete) {
		return config.ErrAccessDenied
	}

	return ss.pluginsystem.DeletePostByID(ID)
}
