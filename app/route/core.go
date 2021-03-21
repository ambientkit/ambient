package route

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/josephspurrier/ambient/app/lib/ambsystem"
	"github.com/josephspurrier/ambient/app/lib/datastorage"
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/lib/websession"
	"github.com/josephspurrier/ambient/assets"
)

// Core -
type Core struct {
	Router  *router.Mux
	Storage *datastorage.Storage
	Render  *htmltemplate.Engine
	Sess    *websession.Session
	Plugins *ambsystem.PluginSystem
}

// Register all routes.
func Register(storage *datastorage.Storage, sess *websession.Session, tmpl *htmltemplate.Engine, mux *router.Mux, plugins *ambsystem.PluginSystem) (*Core, error) {
	// Create core app.
	c := &Core{
		Router:  mux,
		Storage: storage,
		Render:  tmpl,
		Sess:    sess,
		Plugins: plugins,
	}

	// Register routes.
	registerHomePost(&HomePost{c})
	registerStyles(&Styles{c})
	registerAuthUtil(&AuthUtil{c})
	registerXMLUtil(&XMLUtil{c})
	registerAdminPost(&AdminPost{c})
	registerPluginPage(&PluginPage{c})

	// This should be last because it catches all other pages at the root.
	registerPost(&Post{c})

	return c, nil
}

func SetupRouter(tmpl *htmltemplate.Engine) *router.Mux {
	// Set the handling of all responses.
	customServeHTTP := func(w http.ResponseWriter, r *http.Request, status int, err error) {
		// Handle only errors.
		if status >= 400 {
			vars := make(map[string]interface{})
			vars["title"] = fmt.Sprint(status)
			errTemplate := "400"
			if status == 404 {
				errTemplate = "404"
			} else {
				if err != nil {
					fmt.Println(err.Error())
				}

			}
			status, err = tmpl.ErrorTemplate(w, r, "base", errTemplate, vars)
			if err != nil {
				if err != nil {
					log.Println(err.Error())
				}
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}

		// Display server errors.
		if status >= 500 {
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	// Send all 404 to the customer handler.
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customServeHTTP(w, r, http.StatusNotFound, nil)
	})

	// Set up the router.
	rr := router.New(customServeHTTP, notFound)

	// Static assets.
	rr.Get("/assets...", func(w http.ResponseWriter, r *http.Request) (status int, err error) {
		// Don't allow directory browsing.
		if strings.HasSuffix(r.URL.Path, "/") {
			return http.StatusNotFound, nil
		}

		// Use the root directory.
		fsys, err := fs.Sub(assets.CSS, ".")
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Get the requested file name.
		fname := strings.TrimPrefix(r.URL.Path, "/assets/")

		// Open the file.
		f, err := fsys.Open(fname)
		if err != nil {
			return http.StatusNotFound, nil
		}
		defer f.Close()

		// Get the file time.
		st, err := f.Stat()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		http.ServeContent(w, r, fname, st.ModTime(), f.(io.ReadSeeker))
		return
	})

	return rr
}
