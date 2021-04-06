package core

// SavePost saves a post.
func (ss *SecureSite) SavePost(ID string, post Post) error {
	grant := "site.post:write"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	ss.storage.Site.Posts[ID] = post

	return ss.storage.Save()
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(onlyPublished bool) (PostWithIDList, error) {
	grant := "site.post:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PostsAndPages(onlyPublished), nil
}

// PublishedPosts returns the list of published posts.
func (ss *SecureSite) PublishedPosts() ([]Post, error) {
	grant := "site.post:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PublishedPosts(), nil
}

// PublishedPages returns the list of published pages.
func (ss *SecureSite) PublishedPages() ([]Post, error) {
	grant := "site.post:read"

	if !ss.Authorized(grant) {
		return nil, ErrAccessDenied
	}

	return ss.storage.Site.PublishedPages(), nil
}

// PostBySlug returns the post by slug.
func (ss *SecureSite) PostBySlug(slug string) (PostWithID, error) {
	grant := "site.post:read"

	if !ss.Authorized(grant) {
		return PostWithID{}, ErrAccessDenied
	}

	return ss.storage.Site.PostBySlug(slug), nil
}

// PostByID returns the post by ID.
func (ss *SecureSite) PostByID(ID string) (Post, error) {
	grant := "site.post:read"

	if !ss.Authorized(grant) {
		return Post{}, ErrAccessDenied
	}

	post, ok := ss.storage.Site.Posts[ID]
	if !ok {
		return Post{}, ErrNotFound
	}

	return post, nil
}

// DeletePostByID deletes a post.
func (ss *SecureSite) DeletePostByID(ID string) error {
	grant := "site.post:delete"

	if !ss.Authorized(grant) {
		return ErrAccessDenied
	}

	delete(ss.storage.Site.Posts, ID)

	return ss.storage.Save()
}
