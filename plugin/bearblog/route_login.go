package bearblog

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/josephspurrier/ambient/plugin/bearblog/lib/passhash"
	"github.com/josephspurrier/ambient/plugin/bearblog/lib/totp"
)

// login allows a user to login to the dashboard.
func (p *Plugin) login(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := p.Mux.Param(r, "slug")
	loginURL, err := p.Site.PluginField(LoginURL)
	if err != nil {
		return p.Site.Error(err)
	}

	if slug != loginURL {
		return http.StatusNotFound, nil
	}

	vars := make(map[string]interface{})
	vars["title"] = "Login"
	vars["token"] = p.Security.SetCSRF(r)

	return p.Render.PluginPage(w, r, assets, "template/content/login", p.FuncMap(r), vars)
}

func (p *Plugin) loginPost(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := p.Mux.Param(r, "slug")
	loginURL, err := p.Site.PluginField(LoginURL)
	if err != nil {
		return p.Site.Error(err)
	}

	if slug != loginURL {
		return http.StatusNotFound, nil
	}

	r.ParseForm()

	// CSRF protection.
	success := p.Security.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	mfa := r.FormValue("mfa")
	remember := r.FormValue("remember")

	allowedUsername := os.Getenv("AMB_USERNAME")
	if len(allowedUsername) == 0 {
		log.Println("Environment variable missing:", "AMB_USERNAME")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	hash := os.Getenv("AMB_PASSWORD_HASH")
	if len(hash) == 0 {
		log.Println("Environment variable missing:", "AMB_PASSWORD_HASH")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Get the MFA key - if the environment variable doesn't exist, then
	// let the MFA pass.
	mfakey := os.Getenv("AMB_MFA_KEY")
	mfaSuccess := true
	if len(mfakey) > 0 {
		imfa := 0
		imfa, err = strconv.Atoi(mfa)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		mfaSuccess, err = totp.Authenticate(imfa, mfakey)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	// When running locally, let any MFA pass.
	// if envdetect.RunningLocalDev() {
	// 	mfaSuccess = true
	// }

	// Decode the hash - this is to allow it to be stored easily since dollar
	// signs are difficult to work with.
	hashDecoded, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	passMatch := passhash.MatchString(string(hashDecoded), password)

	// If the username and password don't match, then just redirect.
	if username != allowedUsername || !passMatch || !mfaSuccess {
		fmt.Printf("Login attempt failed. Username: %v (expected: %v) | Password match: %v | MFA success: %v\n", username, allowedUsername, passMatch, mfaSuccess)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Printf("Login attempt successful.\n")

	p.Security.SetUser(r, username)
	if remember == "on" {
		p.Security.RememberMe(r, true)
	} else {
		p.Security.RememberMe(r, false)
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

func (p *Plugin) logout(w http.ResponseWriter, r *http.Request) (status int, err error) {
	p.Security.Logout(r)

	http.Redirect(w, r, "/", http.StatusFound)
	return
}
