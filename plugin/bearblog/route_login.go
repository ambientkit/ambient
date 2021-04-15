package bearblog

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/josephspurrier/ambient/plugin/bearblog/lib/passhash"
	"github.com/josephspurrier/ambient/plugin/bearblog/lib/totp"
)

// login allows a user to login to the dashboard.
func (p *Plugin) login(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := p.Mux.Param(r, "slug")
	loginURL, err := p.Site.PluginSetting(LoginURL)
	if err != nil {
		return p.Site.Error(err)
	}

	if slug != loginURL {
		return http.StatusNotFound, nil
	}

	vars := make(map[string]interface{})
	vars["title"] = "Login"
	vars["token"] = p.Site.SetCSRF(r)

	return p.Render.Page(w, r, assets, "template/content/login", p.funcMap(r), vars)
}

func (p *Plugin) loginPost(w http.ResponseWriter, r *http.Request) (status int, err error) {
	slug := p.Mux.Param(r, "slug")
	loginURL, err := p.Site.PluginSetting(LoginURL)
	if err != nil {
		return p.Site.Error(err)
	}

	if slug != loginURL {
		return http.StatusNotFound, nil
	}

	r.ParseForm()

	// CSRF protection.
	success := p.Site.CSRF(r)
	if !success {
		return http.StatusBadRequest, nil
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	mfa := r.FormValue("mfa")
	remember := r.FormValue("remember")

	allowedUsername, err := p.Site.PluginSettingString(Username)
	if err != nil {
		p.Site.Error(err)
		return
	}

	allowedPassword, err := p.Site.PluginSettingString(Password)
	if err != nil {
		p.Site.Error(err)
		return
	}

	mfakey, err := p.Site.PluginSettingString(MFAKey)
	if err != nil {
		p.Site.Error(err)
		return
	}

	// Get the MFA key - if the environment variable doesn't exist, then
	// let the MFA pass.
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

	// Decode the hash - this is to allow it to be stored easily since dollar
	// signs are difficult to work with.
	hashDecoded, err := base64.StdEncoding.DecodeString(allowedPassword)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	passMatch := passhash.MatchString(string(hashDecoded), password)

	// If the username and password don't match, then just redirect.
	if username != allowedUsername || !passMatch || !mfaSuccess {
		p.Log.Info("login attempt failed. Username: %v (expected: %v) | Password match: %v | MFA success: %v", username, allowedUsername, passMatch, mfaSuccess)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = p.Site.UserLogin(r, username)
	if err != nil {
		p.Log.Info("login attempt failed for '%v': %v", username, err.Error())
	} else {
		p.Log.Info("login attempt successful for user: %v", username)
	}
	if remember == "on" {
		err = p.Site.UserPersist(r, true)
	} else {
		err = p.Site.UserPersist(r, false)
	}

	if err != nil {
		p.Log.Info("login persist failed for user '%v': %v", username, err.Error())
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
	return
}

func (p *Plugin) logout(w http.ResponseWriter, r *http.Request) (status int, err error) {
	err = p.Site.UserLogout(r)
	if err != nil {
		p.Log.Info("logout failed: %v", err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return
}
