package hello

import (
	"errors"
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
		post.Title == returnedPost.Title {
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
		post.Title == returnedPost.Title {
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
		post.Title == returnedPost.Title {
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
		post.Title == returnedPost.Title {
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
		post.Title == returnedPost.Title {
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

func (p *Plugin) userPersist(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.UserLogin(r, "username")
	if err != nil {
		return p.Mux.StatusError(http.StatusInternalServerError, err)
	}

	err = p.Site.UserPersist(r, true)
	if err != nil {
		return p.Mux.StatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (p *Plugin) userPersistFalse(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.UserPersist(r, false)
	if err != nil {
		return p.Mux.StatusError(http.StatusInternalServerError, err)
	}

	return nil
}

func (p *Plugin) grantRequests(w http.ResponseWriter, r *http.Request) error {
	requests := p.GrantRequests()
	fmt.Fprintf(w, "Grant requests: %v", len(requests))
	return nil
}

func (p *Plugin) userLogout(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.UserLogin(r, "username")
	if err != nil {
		return p.Mux.StatusError(http.StatusInternalServerError, err)
	}

	err = p.Site.UserLogout(r)
	if err != nil {
		return p.Mux.StatusError(http.StatusInternalServerError, err)
	}

	rAuth, _ := p.Site.AuthenticatedUser(r)
	if rAuth != "" {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("username should not be found:"+rAuth))
	}

	fmt.Fprint(w, "User cleared.")

	return nil
}

func (p *Plugin) logoutAllUsers(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.LogoutAllUsers(r)
	if err != nil {
		return p.Mux.StatusError(http.StatusInternalServerError, err)
	}

	fmt.Fprint(w, "Users cleared.")

	return nil
}

func (p *Plugin) setCSRF(w http.ResponseWriter, r *http.Request) error {
	token := p.Site.SetCSRF(r)
	if len(token) == 0 {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("token is missing"))
	}

	fmt.Fprint(w, token)

	return nil
}

func (p *Plugin) cSRF(w http.ResponseWriter, r *http.Request) error {
	token := r.FormValue("token")
	valid := p.Site.CSRF(r, token)
	if !valid {
		return p.Mux.StatusError(http.StatusBadRequest, errors.New("token is not valid"))
	}

	fmt.Fprintf(w, "Token is valid.")

	return nil
}

func (p *Plugin) sessionValue(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetSessionValue(r, "foo", "bar")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, errors.New("could not set session value"))
	}

	val := p.Site.SessionValue(r, "foo")
	if val != "bar" {
		return p.Mux.StatusError(http.StatusBadRequest, errors.New("could not get session value"))
	}

	p.Site.DeleteSessionValue(r, "foo")

	val = p.Site.SessionValue(r, "foo")
	if val == "bar" {
		return p.Mux.StatusError(http.StatusBadRequest, errors.New("could not delete session value"))
	}

	fmt.Fprint(w, "Session value works.")

	return nil
}

func (p *Plugin) pluginNeighborSettingsList(w http.ResponseWriter, r *http.Request) error {
	settings, err := p.Site.PluginNeighborSettingsList("neighbor")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Neighbor settings: %v", len(settings))

	return nil
}

func (p *Plugin) setPluginSetting(w http.ResponseWriter, r *http.Request) error {
	// Set setting value.
	err := p.Site.SetPluginSetting(Username, "foo")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	// Get string setting.
	val, err := p.Site.PluginSettingString(Username)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}
	if val != "foo" {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("missing string value"))
	}

	// Get bool setting.
	b, err := p.Site.PluginSettingBool(SafeMode)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}
	if !b {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("missing bool false default"))
	}

	// Set setting value.
	err = p.Site.SetPluginSetting(SafeMode, "false")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	// Get bool setting.
	b, err = p.Site.PluginSettingBool(SafeMode)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}
	if b {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("missing bool true value"))
	}

	ival, err := p.Site.PluginSetting(SafeMode)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}
	if ival != "false" {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("missing interface false value"))
	}

	fmt.Fprint(w, "Plugin setting works.")

	return nil
}

func (p *Plugin) setNeighborPluginSetting(w http.ResponseWriter, r *http.Request) error {
	// Set setting value.
	err := p.Site.SetNeighborPluginSetting("neighbor", Username, "foo")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	// Get string setting.
	val, err := p.Site.NeighborPluginSettingString("neighbor", Username)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}
	if val != "foo" {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("missing string value"))
	}

	// Set setting value.
	err = p.Site.SetNeighborPluginSetting("neighbor", SafeMode, "false")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	ival, err := p.Site.NeighborPluginSetting("neighbor", SafeMode)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}
	if ival != "false" {
		return p.Mux.StatusError(http.StatusInternalServerError, errors.New("missing interface false value"))
	}

	fmt.Fprint(w, "Plugin neighbor setting works.")

	return nil
}

func (p *Plugin) pluginTrusted(w http.ResponseWriter, r *http.Request) error {
	trusted, err := p.Site.PluginTrusted("neighbor")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	trusted2, err := p.Site.PluginTrusted("trust")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Plugin trusted: %v %v", trusted, trusted2)

	return nil
}

func (p *Plugin) title(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetTitle("foo")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	title, err := p.Site.Title()
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Site title: %v", title)

	return nil
}

func (p *Plugin) scheme(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetScheme("https")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	scheme, err := p.Site.Scheme()
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Site scheme: %v", scheme)

	return nil
}

func (p *Plugin) uRL(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetURL("bar")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	URL, err := p.Site.URL()
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	FullURL, err := p.Site.FullURL()
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Site URL: %v | Full URL: %v", URL, FullURL)

	return nil
}

func (p *Plugin) updated(w http.ResponseWriter, r *http.Request) error {
	timestamp, err := p.Site.Updated()
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Site updated: %v", timestamp.Format("20060102"))

	return nil
}

func (p *Plugin) content(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetContent("foo bar")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	content, err := p.Site.Content()
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Site content: %v", content)

	return nil
}

func (p *Plugin) tags(w http.ResponseWriter, r *http.Request) error {
	err := p.Site.SetContent("foo bar")
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	tags, err := p.Site.Tags(false)
	if err != nil {
		return p.Mux.StatusError(http.StatusBadRequest, err)
	}

	fmt.Fprintf(w, "Site tags: %v %v", tags[0].Name, tags[0].Timestamp.Format("20060102"))

	return nil
}

func (p *Plugin) assets(w http.ResponseWriter, r *http.Request) error {
	assets, _ := p.Assets()

	fmt.Fprintf(w, "Site assets: %#v", assets)

	return nil
}

func (p *Plugin) assetsHello(w http.ResponseWriter, r *http.Request) error {
	content, err := p.Site.Content()
	if err != nil {
		return p.Site.Error(err)
	}

	if content == "" {
		content = "*No content yet.*"
	}

	vars := make(map[string]interface{})
	vars["postcontent"] = content
	return p.Render.Page(w, r, assets, "template/content/home", p.FuncMap(), vars)
}
