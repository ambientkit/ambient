package bearblog

import (
	"net/http"
	"time"

	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/plugin/bearblog/lib/uuid"
)

func (p *Plugin) postAdminIndex(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Posts"

	postsAndPages, err := p.Site.PostsAndPages(false)
	if err != nil {
		return p.Site.Error(err)
	}

	vars["posts"] = postsAndPages

	return p.Render.Page(w, r, assets, "template/content/bloglist_edit", p.funcMap(r), vars)
}

func (p *Plugin) postAdminCreate(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "New post"
	vars["token"] = p.Site.SetCSRF(r)

	return p.Render.Page(w, r, assets, "template/content/post_create", p.funcMap(r), vars)
}

func (p *Plugin) postAdminStore(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID, err := uuid.Generate()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r.ParseForm()

	// CSRF protection.
	success := p.Site.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	now := time.Now()

	var post ambient.Post
	post.Title = r.FormValue("title")
	post.URL = r.FormValue("slug")
	post.Canonical = r.FormValue("canonical_url")
	post.Created = now
	post.Updated = now
	pubDate := r.FormValue("published_date")
	if pubDate == "" {
		pubDate = now.Format("2006-01-02")
	}
	ts, err := time.Parse("2006-01-02", pubDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	post.Timestamp = ts
	post.Content = r.FormValue("content")
	post.Tags = post.Tags.Split(r.FormValue("tags"))
	post.Page = r.FormValue("is_page") == "on"
	post.Published = r.FormValue("publish") == "on"

	// Save to storage.
	err = p.Site.SavePost(ID, post)
	if err != nil {
		return p.Site.Error(err)
	}

	http.Redirect(w, r, "/dashboard/posts/"+ID, http.StatusFound)
	return
}

func (p *Plugin) postAdminEdit(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Edit post"
	vars["token"] = p.Site.SetCSRF(r)

	ID := p.Mux.Param(r, "id")

	post, err := p.Site.PostByID(ID)
	if err != nil {
		return p.Site.Error(err)
	}

	vars["id"] = ID
	vars["ptitle"] = post.Title
	vars["url"] = post.URL
	vars["canonical"] = post.Canonical
	vars["timestamp"] = post.Timestamp
	vars["body"] = post.Content
	vars["tags"] = post.Tags.String()
	vars["page"] = post.Page
	vars["published"] = post.Published

	return p.Render.Page(w, r, assets, "template/content/post_edit", p.funcMap(r), vars)
}

func (p *Plugin) postAdminUpdate(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID := p.Mux.Param(r, "id")

	post, err := p.Site.PostByID(ID)
	if err != nil {
		return p.Site.Error(err)
	}

	// Save the site.
	r.ParseForm()

	// CSRF protection.
	success := p.Site.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	now := time.Now()

	post.Title = r.FormValue("title")
	post.URL = r.FormValue("slug")
	post.Canonical = r.FormValue("canonical_url")
	post.Updated = now
	pubDate := r.FormValue("published_date")
	ts, err := time.Parse("2006-01-02", pubDate)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	post.Timestamp = ts
	post.Content = r.FormValue("content")
	post.Tags = post.Tags.Split(r.FormValue("tags"))
	post.Page = r.FormValue("is_page") == "on"
	post.Published = r.FormValue("publish") == "on"

	// Save to storage.
	err = p.Site.SavePost(ID, post)
	if err != nil {
		return p.Site.Error(err)
	}

	http.Redirect(w, r, "/dashboard/posts/"+ID, http.StatusFound)
	return
}

func (p *Plugin) postAdminDestroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID := p.Mux.Param(r, "id")

	_, err = p.Site.PostByID(ID)
	if err != nil {
		return p.Site.Error(err)
	}

	err = p.Site.DeletePostByID(ID)
	if err != nil {
		return p.Site.Error(err)
	}

	http.Redirect(w, r, "/dashboard/posts", http.StatusFound)
	return
}
