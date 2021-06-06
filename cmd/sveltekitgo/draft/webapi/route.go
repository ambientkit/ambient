package webapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Routes gets routes for the plugin.
func (p *Plugin) Routes() {
	p.Mux.Get("/", p.index)
	p.Mux.Get("/v1/auth/session", p.session)
	p.Mux.Post("/v1/auth/login", p.login)
}

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
		return http.StatusBadRequest, nil
	}

	if err != nil {
		return http.StatusBadRequest, nil
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

// User -
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// LoadUsers -
func LoadUsers() ([]User, error) {
	b, err := ioutil.ReadFile("cmd/sveltekitgo/storage/users.json")
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	err = json.Unmarshal(b, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func passwordMatch(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
