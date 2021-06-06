package webapi

import (
	"encoding/json"
	"net/http"
)

func (p *Plugin) loadNotes(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type response struct {
		Notes []Note `json:"notes"`
	}

	userEmail, err := p.Site.AuthenticatedUser(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	notes, err := LoadNotes(userEmail)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return p.JSON(w, http.StatusOK, response{
		Notes: notes,
	})
}

func (p *Plugin) createNote(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type request struct {
		Message string `json:"message"`
	}
	type response struct{}

	userEmail, err := p.Site.AuthenticatedUser(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req := new(request)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = CreateNote(userEmail, req.Message)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return p.JSON(w, http.StatusCreated, response{})
}

func (p *Plugin) updateNote(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type request struct {
		Message string `json:"message"`
	}
	type response struct{}

	userEmail, err := p.Site.AuthenticatedUser(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	req := new(request)
	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = UpdateNote(userEmail, p.Mux.Param(r, "noteid"), req.Message)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return p.JSON(w, http.StatusOK, response{})
}

func (p *Plugin) deleteNote(w http.ResponseWriter, r *http.Request) (status int, err error) {
	type response struct{}

	userEmail, err := p.Site.AuthenticatedUser(r)
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = DeleteNote(userEmail, p.Mux.Param(r, "noteid"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	return p.JSON(w, http.StatusOK, response{})
}
