package app

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	// PortHTTP is the HTTP port.
	PortHTTP = ":80"
	// PortHTTPS is the HTTPS port.
	PortHTTPS = ":443"
	// TargetUI is the UI connection info.
	TargetUI = "http://localhost:8080"
	// TargetAPI is the API connection info.
	TargetAPI = "http://localhost:8081"
)

// Proxy is a proxy that supports UI and API.
type Proxy struct {
	targetUI  *url.URL
	targetAPI *url.URL
	proxyUI   http.Handler
	proxyAPI  http.Handler
}

func runProxy() {
	// Send all requests to the proxy handler.
	log.Fatalln(http.ListenAndServe(PortHTTP, NewProxy(TargetUI, TargetAPI)))
}

// NewProxy returns a reverse proxy that allows serving both the API and the UI
// from the same port based on URL.
func NewProxy(targetUI, targetAPI string) *Proxy {
	urlUI, err := url.Parse(targetUI)
	if err != nil {
		log.Fatalln("ui target error:", err)
	}

	urlAPI, err := url.Parse(targetAPI)
	if err != nil {
		log.Fatalln("api target Error:", err)
	}

	// Changed the ErrorHandler to suppress messages like these:
	// 2021/06/06 18:28:33 http: proxy error: context canceled
	// 2021/06/06 18:28:20 http: proxy error: dial tcp [::1]:8080: connect: connection refused

	pui := httputil.NewSingleHostReverseProxy(urlUI)
	pui.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	papi := httputil.NewSingleHostReverseProxy(urlAPI)
	papi.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// log.Printf("http: proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	return &Proxy{
		targetUI:  urlUI,
		targetAPI: urlAPI,
		proxyUI:   pui,
		proxyAPI:  papi,
	}
}

// ServeHTTP handles the requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If the path starts with /api, serve the API.
	if strings.HasPrefix(r.URL.Path, "/api") {
		p.proxyAPI.ServeHTTP(w, r)
		return
	}

	p.proxyUI.ServeHTTP(w, r)
}
