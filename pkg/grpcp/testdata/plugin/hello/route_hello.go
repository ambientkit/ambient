package hello

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/internal/amberror"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "hello world")
	return nil
}

func (p *Plugin) another(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "hello world - another")
	return nil
}

func (p *Plugin) name(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "hello: %v", p.Mux.Param(r, "name"))
	return nil
}

func (p *Plugin) nameOld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello: %v", p.Mux.Param(r, "name"))
}

func (p *Plugin) errorFunc(w http.ResponseWriter, r *http.Request) error {
	return p.Mux.StatusError(http.StatusForbidden, nil)
}

func (p *Plugin) created(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "created: %v", p.Mux.Param(r, "name"))
	return nil
}

func (p *Plugin) headers(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "headers: %#v", len(r.Header))
	return nil
}

func (p *Plugin) formPOST(w http.ResponseWriter, r *http.Request) error {
	body, _ := io.ReadAll(r.Body)
	fmt.Fprintf(w, "body: %#v", string(body))
	return nil
}

func (p *Plugin) formGet(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, html)
	return nil
}

var html = `
<!DOCTYPE html>
<html lang="en">
<head></head>
<body>
	<form method="post">
	<label for="fname">First name:</label>
	<input type="text" id="fname" name="fname" value="a"><br><br>
	<label for="lname">Last name:</label>
	<input type="text" id="lname" name="lname" value="b"><br><br>
	<input type="submit" value="Submit">
	</form>
</body>
</html>
`

func (p *Plugin) login(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.UserLogin(r, "username")
	s, err2 := p.Site.AuthenticatedUser(r)
	fmt.Fprintf(w, "login: (%v) (%v) (%v)", err, s, err2)
	return nil
}

func (p *Plugin) loggedin(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.AuthenticatedUser(r)
	fmt.Fprintf(w, "login: (%v) (%v)", s, err)
	return nil
}

func (p *Plugin) errorsFunc(w http.ResponseWriter, r *http.Request) error {
	errTest := amberror.ErrGrantNotRequested
	err := p.Site.Error(errTest)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "errors: (%v)", "done")
	return nil
}

func (p *Plugin) neighborPluginGrantList(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginGrantList("neighbor")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Grants: %v", len(s))
	return nil
}

func (p *Plugin) neighborPluginGrantListBad(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginGrantList("neighborbad")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Grants: %v", len(s))
	return nil
}

func (p *Plugin) neighborPluginGrants(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginGrants("neighbor")
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Grants: %v", len(s))
	return nil
}

func (p *Plugin) neighborPluginGranted(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginGranted("neighbor", ambient.GrantRouterRouteWrite)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Granted: %v", s)
	return nil
}

func (p *Plugin) neighborPluginGrantedBad(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginGranted("neighbor", ambient.GrantPluginNeighborGrantRead)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Granted: %v", s)
	return nil
}

func (p *Plugin) setNeighborPluginGrantFalse(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetNeighborPluginGrant("neighbor", ambient.GrantRouterRouteWrite, false)
	if err != nil {
		return err
	}

	s, err := p.Site.NeighborPluginGranted("neighbor", ambient.GrantRouterRouteWrite)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Granted: %v", s)
	return nil
}

func (p *Plugin) setNeighborPluginGrantTrue(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetNeighborPluginGrant("neighbor", ambient.GrantRouterRouteWrite, true)
	if err != nil {
		return err
	}
	s, err := p.Site.NeighborPluginGranted("neighbor", ambient.GrantRouterRouteWrite)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Granted: %v", s)
	return nil
}

func (p *Plugin) neighborPluginRequestedGrant(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginRequestedGrant("neighbor", ambient.GrantRouterRouteWrite)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Requested: %v", s)
	return nil
}

func (p *Plugin) neighborPluginRequestedGrantBad(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.NeighborPluginRequestedGrant("neighbor", ambient.GrantPluginNeighborGrantRead)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Requested: %v", s)
	return nil
}

func (p *Plugin) plugins(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.Plugins()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Plugins: %v", len(s))
	return nil
}

func (p *Plugin) pluginNames(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.PluginNames()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Plugin names: %v", len(s))
	return nil
}

func (p *Plugin) deletePlugin(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.DeletePlugin("neighbor")
	fmt.Fprintf(w, "Delete plugin: %v", err)
	return nil
}

func (p *Plugin) deletePluginBad(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.DeletePlugin("neighborBad")
	fmt.Fprintf(w, "Delete plugin: %v", err)
	return nil
}

func (p *Plugin) enablePlugin(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.EnablePlugin("neighbor", true)
	fmt.Fprintf(w, "Enable plugin: %v", err)
	return nil
}

func (p *Plugin) enablePluginBad(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.EnablePlugin("neighborBad", true)
	fmt.Fprintf(w, "Enable plugin: %v", err)
	return nil
}

func (p *Plugin) loadAllPluginPages(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.LoadAllPluginPages()
	fmt.Fprintf(w, "Load pages: %v", err)
	return nil
}

func (p *Plugin) disablePlugin(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.DisablePlugin("neighbor", true)
	fmt.Fprintf(w, "Disable plugin: %v", err)
	return nil
}

func (p *Plugin) disablePluginBad(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.DisablePlugin("neighborBad", true)
	fmt.Fprintf(w, "Disable plugin: %v", err)
	return nil
}

func (p *Plugin) savePost(w http.ResponseWriter, r *http.Request) error {
	post := ambient.Post{
		Title:     "title",
		URL:       "url",
		Canonical: "canonical",
		Created:   time.Now().Truncate(0),
		Updated:   time.Now().Truncate(0),
		Timestamp: time.Now().Truncate(0),
		Content:   "content",
		Published: true,
		Page:      false,
		Tags: ambient.TagList{
			{Name: "tag1", Timestamp: time.Now().Truncate(0)},
		},
	}

	err := p.Site.SavePost("abc", post)
	if err != nil {
		return err
	}

	arr, err := p.Site.PostsAndPages(true)
	if err != nil {
		return err
	}

	returnedPost := arr[0].Post
	if post.Canonical == returnedPost.Canonical &&
		post.Content == returnedPost.Content &&
		post.Title == post.Title {
		fmt.Fprint(w, "Posts are the same.")
	} else {
		fmt.Fprintf(w, "Posts are different (Len: %v): Sent:\n%v\n|\nReceived:\n%v\n", len(arr), post, returnedPost)
	}

	return nil
}

func (p *Plugin) publishedPosts(w http.ResponseWriter, r *http.Request) error {
	post := ambient.Post{
		Title:     "title",
		URL:       "url",
		Canonical: "canonical",
		Created:   time.Now().Truncate(0),
		Updated:   time.Now().Truncate(0),
		Timestamp: time.Now().Truncate(0),
		Content:   "content",
		Published: true,
		Page:      false,
		Tags: ambient.TagList{
			{Name: "tag1", Timestamp: time.Now().Truncate(0)},
		},
	}

	arr, err := p.Site.PublishedPosts()
	if err != nil {
		return err
	}

	returnedPost := arr[0]
	if post.Canonical == returnedPost.Canonical &&
		post.Content == returnedPost.Content &&
		post.Title == post.Title {
		fmt.Fprint(w, "Posts are the same.")
	} else {
		fmt.Fprintf(w, "Posts are different (Len: %v): Sent:\n%v\n|\nReceived:\n%v\n", len(arr), post, returnedPost)
	}

	return nil
}

func (p *Plugin) publishedPages(w http.ResponseWriter, r *http.Request) error {
	post := ambient.Post{
		Title:     "title2",
		URL:       "url2",
		Canonical: "canonical2",
		Created:   time.Now().Truncate(0),
		Updated:   time.Now().Truncate(0),
		Timestamp: time.Now().Truncate(0),
		Content:   "content2",
		Published: true,
		Page:      true,
		Tags: ambient.TagList{
			{Name: "tag1", Timestamp: time.Now().Truncate(0)},
		},
	}

	err := p.Site.SavePost("abc2", post)
	if err != nil {
		return err
	}

	arr, err := p.Site.PublishedPages()
	if err != nil {
		return err
	}

	returnedPost := arr[0]
	if post.Canonical == returnedPost.Canonical &&
		post.Content == returnedPost.Content &&
		post.Title == post.Title {
		fmt.Fprint(w, "Pages are the same.")
	} else {
		fmt.Fprintf(w, "Pages are different (Len: %v): Sent:\n%v\n|\nReceived:\n%v\n", len(arr), post, returnedPost)
	}

	return nil
}

func (p *Plugin) postBySlug(w http.ResponseWriter, r *http.Request) error {
	post := ambient.Post{
		Title:     "title",
		URL:       "url",
		Canonical: "canonical",
		Created:   time.Now().Truncate(0),
		Updated:   time.Now().Truncate(0),
		Timestamp: time.Now().Truncate(0),
		Content:   "content",
		Published: true,
		Page:      false,
		Tags: ambient.TagList{
			{Name: "tag1", Timestamp: time.Now().Truncate(0)},
		},
	}

	returnedPost, err := p.Site.PostBySlug("url")
	if err != nil {
		return err
	}

	if post.Canonical == returnedPost.Canonical &&
		post.Content == returnedPost.Content &&
		post.Title == post.Title {
		fmt.Fprint(w, "Pages are the same.")
	} else {
		fmt.Fprintf(w, "Pages are different: Sent:\n%v\n|\nReceived:\n%v\n", post, returnedPost)
	}

	return nil
}

func (p *Plugin) postBySlugBad(w http.ResponseWriter, r *http.Request) error {
	_, err := p.Site.PostBySlug("urlBad")
	if err != nil {
		return p.Mux.StatusError(http.StatusNotFound, err)
	}

	return nil
}

func (p *Plugin) postByID(w http.ResponseWriter, r *http.Request) error {
	post := ambient.Post{
		Title:     "title",
		URL:       "url",
		Canonical: "canonical",
		Created:   time.Now().Truncate(0),
		Updated:   time.Now().Truncate(0),
		Timestamp: time.Now().Truncate(0),
		Content:   "content",
		Published: true,
		Page:      false,
		Tags: ambient.TagList{
			{Name: "tag1", Timestamp: time.Now().Truncate(0)},
		},
	}

	returnedPost, err := p.Site.PostByID("abc")
	if err != nil {
		return err
	}

	if post.Canonical == returnedPost.Canonical &&
		post.Content == returnedPost.Content &&
		post.Title == post.Title {
		fmt.Fprint(w, "Pages are the same.")
	} else {
		fmt.Fprintf(w, "Pages are different: Sent:\n%v\n|\nReceived:\n%v\n", post, returnedPost)
	}

	return nil
}

func (p *Plugin) postByIDBad(w http.ResponseWriter, r *http.Request) error {
	_, err := p.Site.PostBySlug("abcBad")
	if err != nil {
		return p.Mux.StatusError(http.StatusNotFound, err)
	}

	return nil
}

func (p *Plugin) deletePostByID(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.DeletePostByID("abcBad")
	if err != nil {
		return p.Mux.StatusError(http.StatusNotFound, err)
	}

	returnedPost, err := p.Site.PostByID("abc")
	if err != nil {
		return err
	}

	if returnedPost.Content != "content" {
		return p.Mux.StatusError(http.StatusInternalServerError, fmt.Errorf("post should exist"))
	}

	err = p.Site.DeletePostByID("abc")
	if err != nil {
		return p.Mux.StatusError(http.StatusNotFound, err)
	}

	returnedPost, err = p.Site.PostByID("abc")
	if err != nil {
		fmt.Fprint(w, "Works.")
		return nil
	}

	return p.Mux.StatusError(http.StatusInternalServerError, fmt.Errorf("post should not exist"))
}

func (p *Plugin) pluginNeighborRoutesList(w http.ResponseWriter, r *http.Request) error {
	routes, err := p.Site.PluginNeighborRoutesList("neighbor")
	if err != nil {
		return p.Mux.StatusError(http.StatusNotFound, err)
	}

	fmt.Fprintf(w, "Routes: %v", len(routes))

	return nil
}

func (p *Plugin) pluginNeighborRoutesListBad(w http.ResponseWriter, r *http.Request) error {
	routes, err := p.Site.PluginNeighborRoutesList("neighborBad")
	if err != nil {
		return p.Mux.StatusError(http.StatusNotFound, err)
	}

	fmt.Fprintf(w, "Routes: %v", len(routes))

	return nil
}
