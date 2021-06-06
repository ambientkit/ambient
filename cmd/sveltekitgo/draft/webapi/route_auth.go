package webapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	return p.JSON(w, http.StatusOK, response{
		Status:  http.StatusOK,
		Message: "ok",
	})
}

// Session returns the session information for the current user. The front-end
// uses this to get the status of the server-side cookie.
func (p *Plugin) session(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type response struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		IsLoggedIn   bool   `json:"isLoggedIn"`
		RedirectPath string `json:"redirectPath"`
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
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				IsLoggedIn:   true,
				RedirectPath: "",
			})
		}
	}

	return http.StatusBadRequest, err
}

func (p *Plugin) login(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	req := new(request)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	users, err := LoadUsers()
	if err != nil {
		return http.StatusBadRequest, err
	}

	for _, user := range users {
		if user.Email == req.Email {
			if passwordMatch(req.Password, user.Password) {
				// Set the cookie.
				err = p.Site.UserLogin(r, req.Email)
				if err != nil {
					return http.StatusInternalServerError, err
				}

				return p.JSON(w, http.StatusOK, response{
					FirstName: user.FirstName,
					LastName:  user.LastName,
				})
			}
		}
	}

	return http.StatusBadRequest, fmt.Errorf("user information is not valid")
}

func (p *Plugin) logout(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type response struct {
		Status bool `json:"status"`
	}

	err = p.Site.UserLogout(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return p.JSON(w, http.StatusOK, response{
		Status: true,
	})
}

func (p *Plugin) register(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type request struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password"`
	}
	type response struct {
		Message string `json:"message"`
	}

	req := new(request)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	users, err := LoadUsers()
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Ensure the user doesn't already exist.
	for _, user := range users {
		if user.Email == req.Email {
			return http.StatusBadRequest, errors.New("user already exists")
		}
	}

	err = CreateUser(User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	return p.JSON(w, http.StatusCreated, response{
		Message: "user created",
	})
}
