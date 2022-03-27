package ambient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Toolkit provides utilities to plugins.
type Toolkit struct {
	Log    Logger
	Mux    Router
	Render Renderer
	Site   SecureSite
}

// Redirect to a relative page with the proper URL prefix.
func (t *Toolkit) Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, t.Path(url), code)
}

// Path to a page with the proper URL prefix.
func (t *Toolkit) Path(url string) string {
	// Don't want to use path.Join() because it will strip the trailing slash in
	// some cases.
	return fmt.Sprintf("%v%v", os.Getenv("AMB_URL_PREFIX"), url)
}

// JSON sends a JSON response that is marshalable.
func (t *Toolkit) JSON(w http.ResponseWriter, status int, response interface{}) error {
	// Convert to JSON bytes.
	b, err := json.Marshal(response)
	if err != nil {
		return StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	return t.sendJSON(w, status, b)
}

// JSONPretty sends an indented JSON response that is marshalable.
func (t *Toolkit) JSONPretty(w http.ResponseWriter, status int, response interface{}) error {
	// Convert to JSON bytes.
	b, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	return t.sendJSON(w, status, b)
}

// sendJSON sends a JSON response.
func (t *Toolkit) sendJSON(w http.ResponseWriter, status int, response []byte) error {
	// Set the header.
	w.Header().Set("Content-Type", "application/json")

	// Write out the response.
	_, err := fmt.Fprint(w, string(response))
	if err != nil {
		return StatusError{Code: http.StatusInternalServerError, Err: err}
	}

	return nil
}
