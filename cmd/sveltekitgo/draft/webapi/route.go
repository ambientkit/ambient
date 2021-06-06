package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HelloResponse -
type HelloResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (p *Plugin) index(w http.ResponseWriter, r *http.Request) (status int, err error) {
	resp := &HelloResponse{
		Status:  http.StatusOK,
		Message: "Cool!",
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprint(w, "Error!")
		return
	}

	return fmt.Fprint(w, string(b))
}
