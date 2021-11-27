package ja3rp

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

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

func TestMakeReverseProxyServer(t *testing.T) {
	dsm := &destinationServerMock{}

	// mock destination server
	ds := httptest.NewServer(dsm)
	defer ds.Close()

	// setup reverse proxy server
	s, err := NewReverseProxyServer(ds.URL)
	if err != nil {
		t.Fatal(err)
	}
	addr := "localhost:" + getPort(1337)

	// start listening in the background
	go (func() {
		if err = s.Listen(addr); err != nil {
			t.Fatal(err)
		}
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

func TestMakeServer(t *testing.T) {
	expected := "ok"

	s := NewServer()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	})

	addr := "localhost:" + getPort(1336)

	go (func() {
		if err := s.Listen(addr); err != nil {
			t.Fatal(err)
		}
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
