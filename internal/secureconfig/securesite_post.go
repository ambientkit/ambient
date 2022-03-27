package secureconfig

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/amberror"
)

// SavePost saves a post.
func (ss *SecureSite) SavePost(ID string, post ambient.Post) error {
	if !ss.Authorized(ambient.GrantSitePostWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SavePost(ID, post)
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(onlyPublished bool) (ambient.PostWithIDList, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PostsAndPages(onlyPublished), nil
}

// PublishedPosts returns the list of published posts.
func (ss *SecureSite) PublishedPosts() ([]ambient.Post, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPosts(), nil
}

// PublishedPages returns the list of published pages.
func (ss *SecureSite) PublishedPages() ([]ambient.Post, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPages(), nil
}

// PostBySlug returns the post by slug.
func (ss *SecureSite) PostBySlug(slug string) (ambient.PostWithID, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return ambient.PostWithID{}, amberror.ErrAccessDenied
	}

	post := ss.pluginsystem.PostBySlug(slug)
	if post.ID == "" {
		return ambient.PostWithID{}, amberror.ErrNotFound
	}

	return post, nil
}

// PostByID returns the post by ID.
func (ss *SecureSite) PostByID(ID string) (ambient.Post, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return ambient.Post{}, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PostByID(ID)
}

// DeletePostByID deletes a post.
func (ss *SecureSite) DeletePostByID(ID string) error {
	if !ss.Authorized(ambient.GrantSitePostDelete) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.DeletePostByID(ID)
}
