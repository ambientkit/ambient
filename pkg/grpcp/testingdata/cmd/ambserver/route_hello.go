package main

import (
	"fmt"
	"io"
	"net/http"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) error {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	fmt.Fprint(w, "hello world")
	return nil
	//return p.Render.Page(w, r, assets, "template/hello.tmpl", nil, vars)
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
	//fmt.Fprint(w, "hello error page")
	//p.Mux.Error(http.StatusForbidden, w, r)
	//return nil
	return p.Mux.StatusError(http.StatusForbidden, nil)
}

func (p *Plugin) created(w http.ResponseWriter, r *http.Request) error {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "created: %v", p.Mux.Param(r, "name"))
	return nil
	//return p.Render.Page(w, r, assets, "template/hello.tmpl", nil, vars)
}

func (p *Plugin) headersPOST(w http.ResponseWriter, r *http.Request) error {
	//fmt.Fprintf(w, "headers: %#v", r.Header)
	body, _ := io.ReadAll(r.Body)
	fmt.Fprintf(w, "body: %#v", string(body))
	return nil
}

func (p *Plugin) headers(w http.ResponseWriter, r *http.Request) error {
	//fmt.Fprintf(w, "headers: %#v", r.Header)
	// body, _ := io.ReadAll(r.Body)
	// fmt.Fprintf(w, "body: %#v", body)
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
