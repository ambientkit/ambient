package secureconfig

import (
	"context"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/amberror"
)

// SetTitle sets the title.
func (ss *SecureSite) SetTitle(ctx context.Context, title string) error {
	if !ss.Authorized(ctx, ambient.GrantSiteTitleWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetTitle(title)
}

// Title returns the title.
func (ss *SecureSite) Title(ctx context.Context) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantSiteTitleRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Title(), nil
}

// SetScheme sets the site scheme.
func (ss *SecureSite) SetScheme(ctx context.Context, scheme string) error {
	if !ss.Authorized(ctx, ambient.GrantSiteSchemeWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetScheme(scheme)
}

// Scheme returns the site scheme.
func (ss *SecureSite) Scheme(ctx context.Context) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantSiteSchemeRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Scheme(), nil
}

// SetURL sets the site URL.
func (ss *SecureSite) SetURL(ctx context.Context, URL string) error {
	if !ss.Authorized(ctx, ambient.GrantSiteURLWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetURL(URL)
}

// URL returns the URL without the scheme at the beginning.
func (ss *SecureSite) URL(ctx context.Context) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantSiteURLRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.URL(), nil
}

// FullURL returns the URL with the scheme at the beginning.
func (ss *SecureSite) FullURL(ctx context.Context) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantSiteURLRead) ||
		!ss.Authorized(ctx, ambient.GrantSiteSchemeRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.FullURL(), nil
}

// Updated returns the home last updated timestamp.
func (ss *SecureSite) Updated(ctx context.Context) (time.Time, error) {
	if !ss.Authorized(ctx, ambient.GrantSiteUpdatedRead) {
		return time.Now(), amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Updated(), nil
}

// SetContent sets the home page content.
func (ss *SecureSite) SetContent(ctx context.Context, content string) error {
	if !ss.Authorized(ctx, ambient.GrantSiteContentWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetContent(content)
}

// Content returns the site home page content.
func (ss *SecureSite) Content(ctx context.Context) (string, error) {
	if !ss.Authorized(ctx, ambient.GrantSiteContentRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Content(), nil
}

// Tags returns the list of tags.
func (ss *SecureSite) Tags(ctx context.Context, onlyPublished bool) (ambient.TagList, error) {
	if !ss.Authorized(ctx, ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Tags(onlyPublished), nil
}
