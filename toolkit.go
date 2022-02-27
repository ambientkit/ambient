package ambient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
)

// Toolkit provides utilities to plugins.
type Toolkit struct {
	Log    Logger
	Mux    Router
	Render Renderer
	Site   SecureSite
}

// Redirect to a page with the proper URL prefix.
func (t *Toolkit) Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, t.Path(url), code)
}

// Path to a page with the proper URL prefix.
func (t *Toolkit) Path(url string) string {
	return path.Join(os.Getenv("AMB_URL_PREFIX"), url)
}

// JSON sends a JSON response that is marshalable.
func (t *Toolkit) JSON(w http.ResponseWriter, status int, response interface{}) (int, error) {
	// Convert to JSON bytes.
	b, err := json.Marshal(response)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return t.sendJSON(w, status, b)
}

// JSONPretty sends an indented JSON response that is marshalable.
func (t *Toolkit) JSONPretty(w http.ResponseWriter, status int, response interface{}) (int, error) {
	// Convert to JSON bytes.
	b, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return t.sendJSON(w, status, b)
}

// sendJSON sends a JSON response.
func (t *Toolkit) sendJSON(w http.ResponseWriter, status int, response []byte) (int, error) {
	// Set the header.
	w.Header().Set("Content-Type", "application/json")

	// Write out the response.
	_, err := fmt.Fprint(w, string(response))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Return the status.
	return status, nil
}
