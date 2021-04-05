package bearblog

import (
	"net/http"
	"time"

	"github.com/josephspurrier/ambient/app/lib/uuid"
	"github.com/josephspurrier/ambient/app/model"
)

func (p *Plugin) postAdminIndex(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "Posts"

	postsAndPages, err := p.Site.PostsAndPages(false)
	if err != nil {
		return p.Site.Error(err)
	}

	vars["posts"] = postsAndPages

	return p.Render.PluginPage(w, r, assets, "template/content/bloglist_edit", p.FuncMap(r), vars)
}

func (p *Plugin) postAdminCreate(w http.ResponseWriter, r *http.Request) (status int, err error) {
	vars := make(map[string]interface{})
	vars["title"] = "New post"
	vars["token"] = p.Security.SetCSRF(r)

	return p.Render.PluginPage(w, r, assets, "template/content/post_create", p.FuncMap(r), vars)
}

func (p *Plugin) postAdminStore(w http.ResponseWriter, r *http.Request) (status int, err error) {
	ID, err := uuid.Generate()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r.ParseForm()

	// CSRF protection.
	success := p.Security.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	now := time.Now()

	var post model.Post
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
	vars["token"] = p.Security.SetCSRF(r)

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

	return p.Render.PluginPage(w, r, assets, "templates/content/post_edit", p.FuncMap(r), vars)
}

// func (p *Plugin) postAdminUpdate(w http.ResponseWriter, r *http.Request) (status int, err error) {
// 	ID := p.Router.Param(r, "id")

// 	var p model.Post
// 	var ok bool
// 	if p, ok = p.Storage.Site.Posts[ID]; !ok {
// 		return http.StatusNotFound, nil
// 	}

// 	// Save the site.
// 	r.ParseForm()

// 	// CSRF protection.
// 	success := p.Sess.CSRF(r)
// 	if !success {
// 		return http.StatusBadRequest, nil
// 	}

// 	now := time.Now()

// 	p.Title = r.FormValue("title")
// 	p.URL = r.FormValue("slug")
// 	p.Canonical = r.FormValue("canonical_url")
// 	p.Updated = now
// 	pubDate := r.FormValue("published_date")
// 	ts, err := time.Parse("2006-01-02", pubDate)
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}
// 	p.Timestamp = ts
// 	p.Content = r.FormValue("content")
// 	p.Tags = p.Tags.Split(r.FormValue("tags"))
// 	p.Page = r.FormValue("is_page") == "on"
// 	p.Published = r.FormValue("publish") == "on"

// 	p.Storage.Site.Posts[ID] = p

// 	err = p.Storage.Save()
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	http.Redirect(w, r, "/dashboard/posts/"+ID, http.StatusFound)
// 	return
// }

// func (p *Plugin) postAdminDestroy(w http.ResponseWriter, r *http.Request) (status int, err error) {
// 	ID := p.Mux.Param(r, "id")

// 	var ok bool
// 	if _, ok = p.Storage.Site.Posts[ID]; !ok {
// 		return http.StatusNotFound, nil
// 	}

// 	delete(p.Storage.Site.Posts, ID)

// 	err = p.Storage.Save()
// 	if err != nil {
// 		return http.StatusInternalServerError, err
// 	}

// 	http.Redirect(w, r, "/dashboard/posts", http.StatusFound)
// 	return
// }
