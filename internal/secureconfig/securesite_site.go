package secureconfig

import (
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/amberror"
)

// SetTitle sets the title.
func (ss *SecureSite) SetTitle(title string) error {
	if !ss.Authorized(ambient.GrantSiteTitleWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetTitle(title)
}

// Title returns the title.
func (ss *SecureSite) Title() (string, error) {
	if !ss.Authorized(ambient.GrantSiteTitleRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Title(), nil
}

// SetScheme sets the site scheme.
func (ss *SecureSite) SetScheme(scheme string) error {
	if !ss.Authorized(ambient.GrantSiteSchemeWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetScheme(scheme)
}

// Scheme returns the site scheme.
func (ss *SecureSite) Scheme() (string, error) {
	if !ss.Authorized(ambient.GrantSiteSchemeRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Scheme(), nil
}

// SetURL sets the site URL.
func (ss *SecureSite) SetURL(URL string) error {
	if !ss.Authorized(ambient.GrantSiteURLWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetURL(URL)
}

// URL returns the URL without the scheme at the beginning.
func (ss *SecureSite) URL() (string, error) {
	if !ss.Authorized(ambient.GrantSiteURLRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.URL(), nil
}

// FullURL returns the URL with the scheme at the beginning.
func (ss *SecureSite) FullURL() (string, error) {
	if !ss.Authorized(ambient.GrantSiteURLRead) || !ss.Authorized(ambient.GrantSiteSchemeRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.FullURL(), nil
}

// Updated returns the home last updated timestamp.
func (ss *SecureSite) Updated() (time.Time, error) {
	if !ss.Authorized(ambient.GrantSiteUpdatedRead) {
		return time.Now(), amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Updated(), nil
}

// Tags returns the list of tags.
func (ss *SecureSite) Tags(onlyPublished bool) (ambient.TagList, error) {
	if !ss.Authorized(ambient.GrantSitePostRead) {
		return nil, amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Tags(onlyPublished), nil
}

// SetContent sets the home page content.
func (ss *SecureSite) SetContent(content string) error {
	if !ss.Authorized(ambient.GrantSiteContentWrite) {
		return amberror.ErrAccessDenied
	}

	return ss.pluginsystem.SetContent(content)
}

// Content returns the site home page content.
func (ss *SecureSite) Content() (string, error) {
	if !ss.Authorized(ambient.GrantSiteContentRead) {
		return "", amberror.ErrAccessDenied
	}

	return ss.pluginsystem.Content(), nil
}
