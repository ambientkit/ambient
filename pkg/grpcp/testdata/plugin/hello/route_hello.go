package hello

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ambientkit/ambient/internal/config"
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
	errTest := config.ErrGrantNotRequested
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
