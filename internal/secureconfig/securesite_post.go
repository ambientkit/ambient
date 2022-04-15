package secureconfig

import (
	"context"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/amberror"
)

// SavePost saves a post.
func (ss *SecureSite) SavePost(ctx context.Context, ID string, post ambient.Post) error {
	if !ss.Authorized(ctx, ambient.GrantSitePostWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SavePost(ID, post)
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(ctx context.Context, onlyPublished bool) (ambient.PostWithIDList, error) {
	if !ss.Authorized(ctx, ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PostsAndPages(onlyPublished), nil
}

// PublishedPosts returns the list of published posts.
func (ss *SecureSite) PublishedPosts(ctx context.Context) ([]ambient.Post, error) {
	if !ss.Authorized(ctx, ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPosts(), nil
}

// PublishedPages returns the list of published pages.
func (ss *SecureSite) PublishedPages(ctx context.Context) ([]ambient.Post, error) {
	if !ss.Authorized(ctx, ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPages(), nil
}

// PostBySlug returns the post by slug.
func (ss *SecureSite) PostBySlug(ctx context.Context, slug string) (ambient.PostWithID, error) {
	if !ss.Authorized(ctx, ambient.GrantSitePostRead) {
		return ambient.PostWithID{}, amberror.ErrAccessDenied
	}

	post := ss.pluginsystem.PostBySlug(slug)
	if post.ID == "" {
		return ambient.PostWithID{}, amberror.ErrNotFound
	}

	return post, nil
}

// PostByID returns the post by ID.
func (ss *SecureSite) PostByID(ctx context.Context, ID string) (ambient.Post, error) {
	if !ss.Authorized(ctx, ambient.GrantSitePostRead) {
		return ambient.Post{}, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.PostByID(ID)
}

// DeletePostByID deletes a post.
func (ss *SecureSite) DeletePostByID(ctx context.Context, ID string) error {
	if !ss.Authorized(ctx, ambient.GrantSitePostDelete) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.DeletePostByID(ID)
}
