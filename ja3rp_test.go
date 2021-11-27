package ja3rp

import (
	"fmt"
	"github.com/sleeyax/ja3rp/net/http"
	"github.com/sleeyax/ja3rp/net/http/httptest"
	"io"
	"net/url"
	"os"
	"strconv"
	"testing"
)

var testPort = 1337

type mock struct {
	reached bool
}

type destinationServerMock struct {
	mock
}

// getPort gets an available port by environment variable or uses the given fallback value if it's not set.
func getPort(defaultValue int) string {
	if v, ok := os.LookupEnv("TEST_SERVER_PORT"); ok {
		return v
	}

	return strconv.Itoa(defaultValue)
}

func (dsm *destinationServerMock) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	dsm.reached = true
	fmt.Fprintf(w, "ok")
}

func TestReverseProxyServer(t *testing.T) {
	dsm := &destinationServerMock{}

	// mock destination server
	ds := httptest.NewServer(dsm)
	defer ds.Close()

	addr := "localhost:" + getPort(testPort)

	// setup reverse proxy server
	u, err := url.Parse(ds.URL)
	if err != nil {
		t.Fatal(err)
	}
	s := NewServer(addr, ServerOptions{
		Destination: u,
	})
	defer s.Close()

	// start listening in the background
	go (func() {
		s.ListenAndServe()
	})()

	// send HTTP request
	res, err := http.Get("http://" + addr)
	if err != nil {
		t.Fatal(err)
	}

	// verify HTTP response
	if res.StatusCode != 200 {
		t.Fail()
	}
	if !dsm.reached {
		t.Errorf("destination server was not reached")
	}
}

func TestServer(t *testing.T) {
	expected := "ok"
	addr := "localhost:" + getPort(testPort)

	mux := NewMux()

	s := NewServer(addr, ServerOptions{
		Mux: mux,
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	})

	go (func() {
		s.ListenAndServe()
	})()

	res, err := http.Get("http://" + addr)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fail()
	}

	body, _ := io.ReadAll(res.Body)
	if bodyStr := string(body); bodyStr != expected {
		t.Errorf("Invalid body. Expected '%s' but got '%s'", expected, bodyStr)
	}
}
