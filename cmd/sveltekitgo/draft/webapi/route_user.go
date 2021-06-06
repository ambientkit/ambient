package webapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (p *Plugin) userProfile(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type response struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	userEmail, err := p.Site.AuthenticatedUser(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	users, err := LoadUsers()
	if err != nil {
		return http.StatusBadRequest, err
	}

	for _, user := range users {
		if user.Email == userEmail {
			return p.JSON(w, http.StatusOK, response{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			})
		}
	}

	return http.StatusBadRequest, err
}

func (p *Plugin) updateUserProfile(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type request struct {
		FirstName       string `json:"firstName"`
		LastName        string `json:"lastName"`
		Email           string `json:"email"`
		NewPassword     string `json:"newPassword"`
		CurrentPassword string `json:"currentPassword"`
	}
	type response struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	userEmail, err := p.Site.AuthenticatedUser(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req := new(request)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	cu, err := GetUser(userEmail)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if !passwordMatch(req.CurrentPassword, cu.Password) {
		return http.StatusBadRequest, errors.New("password doesn't match")
	}

	pass := ""
	if len(req.NewPassword) > 0 {
		pass = req.NewPassword
	}

	err = UpdateUser(User{
		ID:        cu.ID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  pass,
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return p.JSON(w, http.StatusOK, response{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
}
