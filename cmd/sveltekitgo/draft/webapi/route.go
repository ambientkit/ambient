package webapi

import (
	"net/http"
)

// HelloResponse -
type HelloResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return p.JSON(w, http.StatusOK, HelloResponse{
		Status:  http.StatusOK,
		Message: "Cool!",
	})
}
