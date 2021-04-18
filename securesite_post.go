package ambient

// SavePost saves a post.
func (ss *SecureSite) SavePost(ID string, post Post) error {
	if !ss.Authorized(GrantSitePostWrite) {
		return ErrAccessDenied
	}

	return ss.pluginsystem.SavePost(ID, post)
}

// PostsAndPages returns the list of posts and pages.
func (ss *SecureSite) PostsAndPages(onlyPublished bool) (PostWithIDList, error) {
	if !ss.Authorized(GrantSitePostRead) {
		return nil, ErrAccessDenied
	}

	return ss.pluginsystem.PostsAndPages(onlyPublished), nil
}

// PublishedPosts returns the list of published posts.
func (ss *SecureSite) PublishedPosts() ([]Post, error) {
	if !ss.Authorized(GrantSitePostRead) {
		return nil, ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPosts(), nil
}

// PublishedPages returns the list of published pages.
func (ss *SecureSite) PublishedPages() ([]Post, error) {
	if !ss.Authorized(GrantSitePostRead) {
		return nil, ErrAccessDenied
	}

	return ss.pluginsystem.PublishedPages(), nil
}

// PostBySlug returns the post by slug.
func (ss *SecureSite) PostBySlug(slug string) (PostWithID, error) {
	if !ss.Authorized(GrantSitePostRead) {
		return PostWithID{}, ErrAccessDenied
	}

	return ss.pluginsystem.PostBySlug(slug), nil
}

// PostByID returns the post by ID.
func (ss *SecureSite) PostByID(ID string) (Post, error) {
	if !ss.Authorized(GrantSitePostRead) {
		return Post{}, ErrAccessDenied
	}

	return ss.pluginsystem.PostByID(ID)
}

// DeletePostByID deletes a post.
func (ss *SecureSite) DeletePostByID(ID string) error {
	if !ss.Authorized(GrantSitePostDelete) {
		return ErrAccessDenied
	}

	return ss.pluginsystem.DeletePostByID(ID)
}
