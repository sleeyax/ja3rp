package ja3rp

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/sleeyax/ja3rp/net/http"
	"github.com/sleeyax/ja3rp/net/http/httputil"
	"net/url"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type ServerOptions struct {
	// Target server to forward valid traffic to.
	// The reverse proxy mode will be disabled if this field is nil.
	Destination *url.URL

	// Custom Mux to use.
	Mux *Mux

	// Whitelisted JA3 hashes.
	// Only traffic that matches a JA3 from this list will be accepted.
	// If both Whitelist and Blacklist are specified, Blacklist will precede.
	// If both Whitelist and Blacklist are unspecified, all traffic will go through.
	Whitelist []string

	// Blacklisted JA3 hashes
	// Traffic that matches a JA3 from this list will be ignored.
	// If both Whitelist and Blacklist are specified, Blacklist will precede.
	// If both Whitelist and Blacklist are unspecified, all traffic will go through.
	Blacklist []string

	// Called when a JA3 is found on the Blacklist or not on the Whitelist.
	OnBlocked Handler
}

func (s ServerOptions) handleRoot(w http.ResponseWriter, r *http.Request) {
	ja3Hash := JA3Digest(r.JA3)

	if inArray(ja3Hash, s.Blacklist) || (len(s.Whitelist) > 0 && !inArray(ja3Hash, s.Whitelist)) {
		w.WriteHeader(http.StatusForbidden)

		if s.OnBlocked != nil {
			s.OnBlocked(w, r)
		} else {
			fmt.Fprintf(w, "Access forbidden.")
		}

		return
	}

	if s.Destination == nil {
		fmt.Fprintf(w, "Access granted. JA3 hash: "+ja3Hash)
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
	}

	options.Mux.HandleFunc("/", options.handleRoot)

	return &http.Server{
		Addr:    addr,
		Handler: options.Mux,
	}
}

// JA3Digest creates the JA3 hash from given plaintext JA3 string.
func JA3Digest(ja3 string) string {
	h := md5.New()
	h.Write([]byte(ja3))
	return hex.EncodeToString(h.Sum(nil))
}
