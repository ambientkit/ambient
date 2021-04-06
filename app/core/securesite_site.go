package core

import "time"

// Title returns the title.
func (ss *SecureSite) Title() (string, error) {
	grant := "site.title:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Title, nil
}

// SetScheme sets the site scheme.
func (ss *SecureSite) SetScheme(scheme string) error {
	grant := "site.scheme:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Scheme = scheme

	return ss.storage.Save()
}

// Scheme returns the site scheme.
func (ss *SecureSite) Scheme() (string, error) {
	grant := "site.scheme:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Scheme, nil
}

// SetTitle sets the title.
func (ss *SecureSite) SetTitle(title string) error {
	grant := "site.title:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Title = title

	return ss.storage.Save()
}

// SetURL sets the site URL.
func (ss *SecureSite) SetURL(URL string) error {
	grant := "site.url:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.URL = URL

	return ss.storage.Save()
}

// URL returns the URL without the scheme at the beginning.
func (ss *SecureSite) URL() (string, error) {
	grant := "site.url:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.URL, nil
}

// FullURL returns the URL with the scheme at the beginning.
func (ss *SecureSite) FullURL() (string, error) {
	grant1 := "site.url:read"
	grant2 := "site.scheme:read"

	if !ss.Authorized(grant1) || !ss.Authorized(grant2) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.SiteURL(), nil
}

// Updated returns the home last updated timestamp.
func (ss *SecureSite) Updated() (time.Time, error) {
	grant := "site.updated:read"

	if !ss.Authorized(grant) {
		return time.Now(), ErrAccessDenied
	}

	return ss.storage.Site.Updated, nil
}

// Tags returns the list of tags.
func (ss *SecureSite) Tags(onlyPublished bool) (TagList, error) {
	grant := "site.tags:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.Tags(onlyPublished), nil
}

// SetContent sets the home page content.
func (ss *SecureSite) SetContent(content string) error {
	grant := "site.content:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Content = content

	return ss.storage.Save()
}

// Content returns the site home page content.
func (ss *SecureSite) Content() (string, error) {
	grant := "site.content:read"

	if !ss.Authorized(grant) {
		return "", ErrAccessDenied
	}

	return ss.storage.Site.Content, nil
}
