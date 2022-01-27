package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Paths struct {
	path   string
	method string
}

var p = []Paths{
	{path: "/", method: "GET"},
	{path: "/jobs", method: "GET"},
	{path: "/jobs", method: "POST"},
}

func checkstatuscodes(t *testing.T, path, method string) {
	r := newRouter()

	mockServer := httptest.NewServer(r)

	switch method {
	case "GET":
		resp, err := http.Get(mockServer.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status should be ok, got %d instead", resp.StatusCode)
		}
	case "POST":
		resp, err := http.Post(mockServer.URL+path, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("status should be ok, got %d instead", resp.StatusCode)
		}

	}
}

func TestRoutes(t *testing.T) {
	for _, i := range p {
		checkstatuscodes(t, i.path, i.method)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// We want to hit the `GET /assets/` route to get the index.html file response
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// It isn't wise to test the entire content of the HTML file.
	// Instead, we test that the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}
