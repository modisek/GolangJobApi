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
	// {path: "/jobs", method: "POST"},
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
