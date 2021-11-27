package ja3rp

import (
	"fmt"
	"github.com/sleeyax/ja3rp/net/http"
	"github.com/sleeyax/ja3rp/net/http/httputil"
	"net/url"
)

type ServerOptions struct {
	// Target server to forward valid traffic to.
	// The reverse proxy mode will be disabled if this field is left empty.
	Destination *url.URL

	// Custom Mux to use.
	Mux *Mux
}

func (s *ServerOptions) handleRoot(w http.ResponseWriter, r *http.Request) {
	if s.Destination == nil {
		fmt.Fprintf(w, "ja3rp is running (https://github.com/sleeyax/ja3rp)")
		return
	}

	r.URL.Host = s.Destination.Host
	r.URL.Scheme = s.Destination.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = s.Destination.Host

	proxy := httputil.NewSingleHostReverseProxy(s.Destination)
	proxy.ServeHTTP(w, r)
}

func NewServer(addr string, options ServerOptions) *http.Server {
	if options.Mux == nil {
		options.Mux = NewMux()
		options.Mux.HandleFunc("/", options.handleRoot)
	}

	return &http.Server{
		Addr:    addr,
		Handler: options.Mux,
	}
}
