package hello

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ambientkit/ambient/internal/config"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) error {
	vars := make(map[string]interface{})
	vars["title"] = "Plugins"
	fmt.Fprint(w, "hello world")
	return nil
	//return p.Render.Page(w, r, assets, "template/hello", nil, vars)
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
	//return p.Render.Page(w, r, assets, "template/hello", nil, vars)
}

func (p *Plugin) headersPOST(w http.ResponseWriter, r *http.Request) error {
	//fmt.Fprintf(w, "headers: %#v", r.Header)
	body, _ := io.ReadAll(r.Body)
	fmt.Fprintf(w, "body: %#v", string(body))
	return nil
}

func (p *Plugin) headers(w http.ResponseWriter, r *http.Request) error {
	// fmt.Fprintf(w, "headers: %#v", r.Header)
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

func (p *Plugin) loggedin(w http.ResponseWriter, r *http.Request) error {
	s, err := p.Site.AuthenticatedUser(r)
	fmt.Fprintf(w, "login: (%v) (%v)", s, err)
	return nil
}

func (p *Plugin) errorsFunc(w http.ResponseWriter, r *http.Request) error {
	/*
	   // Error handles returning the proper error.
	   func (ss *SecureSite) Error(siteError error) (err error) {
	   	switch siteError {
	   	case config.ErrAccessDenied, config.ErrGrantNotRequested, config.ErrSettingNotSpecified:
	   		return ambient.StatusError{Code: http.StatusForbidden, Err: siteError}
	   	case config.ErrNotFound:
	   		return ambient.StatusError{Code: http.StatusNotFound, Err: siteError}
	   	default:
	   		return ambient.StatusError{Code: http.StatusInternalServerError, Err: siteError}
	   	}
	   }
	*/

	errTest := config.ErrGrantNotRequested

	err := p.Site.Error(errTest)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "errors: (%v)", "done")
	return nil
}
