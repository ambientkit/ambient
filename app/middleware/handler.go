// Package middleware provides handlers for routes.
package middleware

import (
	"github.com/josephspurrier/ambient/app/lib/htmltemplate"
	"github.com/josephspurrier/ambient/app/lib/router"
	"github.com/josephspurrier/ambient/app/lib/websession"
)

// Handler represents a middleware handler.
type Handler struct {
	Router     *router.Mux
	Render     *htmltemplate.Engine
	Sess       *websession.Session
	SiteURL    string
	SiteScheme string
}

// NewHandler returns a middleware handler.
func NewHandler(te *htmltemplate.Engine, sess *websession.Session, mux *router.Mux, siteURL string, siteScheme string) *Handler {
	return &Handler{
		Render:     te,
		Router:     mux,
		Sess:       sess,
		SiteURL:    siteURL,
		SiteScheme: siteScheme,
	}

}
