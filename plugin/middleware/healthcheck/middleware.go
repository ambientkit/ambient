package healthcheck

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// swagger:route GET /api/healthcheck healthcheck healthcheckGET
//
// Returns an OK status message.
//
// Responses:
//   200: healthcheckResponse
//   400: errorResponse
//   500: errorResponse
func (p *Plugin) healthcheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/healthcheck" {
			data := new(healthcheckResponse).Body
			data.Message = "ok"
			JSON(w, data)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// swagger:response healthcheckResponse
type healthcheckResponse struct {
	// in: body
	Body struct {
		// Health check.
		//
		// required: true
		// example: ok
		Message string `json:"message"`
	}
}

// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Body struct {
		// Error message.
		//
		// required: true
		// example: an error occurred
		Message string `json:"message"`
	}
}

// JSON sends a JSON response.
func JSON(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(b))
}

// ErrorJSON sends a JSON error response.
func ErrorJSON(w http.ResponseWriter, err error, status int) {
	log.Println(err.Error())

	w.WriteHeader(status)
	data := new(errorResponse).Body
	data.Message = err.Error()
	out, _ := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(out))
}
