package ja3rp

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	mux *Mux

	// Enable reverse proxy mode.
	EnableReverseProxy bool

	// Target server to forward valid traffic to.
	// This field will be ignored if EnableReverseProxy is false.
	Destination *url.URL
}

func (s *Server) registerDefaultHandlers() {
	s.HandleFunc("/", s.handleRoot)
}

// NewReverseProxyServer creates a new HTTP reverse proxy server.
func NewReverseProxyServer(destinationUrl string) (*Server, error) {
	u, err := url.Parse(destinationUrl)
	if err != nil {
		return &Server{}, err
	}

	s := &Server{
		EnableReverseProxy: true,
		Destination:        u,
		mux:                NewMux(),
	}

	s.registerDefaultHandlers()

	return s, nil
}

// NewServer creates a regular HTTP server.
// Use NewReverseProxyServer if you want to enable the reverse proxy.
func NewServer() *Server {
	s := &Server{mux: NewMux()}
	s.registerDefaultHandlers()
	return s
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if !s.EnableReverseProxy {
		fmt.Fprintf(w, "ja3rp is running (https://github.com/sleeyax/ja3rp)")
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(s.Destination)

	r.URL.Host = s.Destination.Host
	r.URL.Scheme = s.Destination.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = s.Destination.Host

	proxy.ServeHTTP(w, r)
}

// HandleFunc registers the handler function for the given pattern
// If a handler was already registered for given pattern, it will be overwritten.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) Listen(address string) error {
	return http.ListenAndServe(address, s.mux)
}
