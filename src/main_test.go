package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var wrongPaths = []string{
	"/.well-known/matrix/somewhere",
	"/.well-known/matrixes/client",
	"/.well-knowna/matrixes/client",
	"/well-known/matrix/client",
	"/wellknown/matrix/client",
	"/well-known/matrix/clients",
	"/.well-known/matrix/servers",
	"/randomhere",
	"/",
}

var clientPath = "/.well-known/matrix/client"

func TestWrongPath(t *testing.T) {
	for _, test := range wrongPaths {
		req := httptest.NewRequest(http.MethodGet, test, nil)
		w := httptest.NewRecorder()
		requestHandler(w, req)
		result := w.Result()
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("Wrong status code at wrong path at path %s with status code %d", test, result.StatusCode)
		}
	}
}

func TestWithIdentityServer(t *testing.T) {
	os.Setenv("CLIENT_HOMESERVER", "rando")
	os.Setenv("CLIENT_IDENTITYSERVER", "rando2")

	req := httptest.NewRequest(http.MethodGet, clientPath, nil)
	w := httptest.NewRecorder()
	requestHandler(w, req)
	result := w.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code with status code %d that should be 200", result.StatusCode)
	}
	if header := result.Header.Get("Access-Control-Allow-Origin"); header != "*" {
		t.Errorf("Wrong header received at client")
	}
	if header := result.Header.Get("content-type"); header != "application/json;charset=UTF-8" {
		t.Errorf("Wrong header received at client")
	}

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}

	var clientJSON ClientResponse
	err = json.Unmarshal(data, &clientJSON)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if clientJSON.Homeserver.BaseURL != "rando" {
		t.Errorf("Wrong homeserver received")
	}
	if clientJSON.IdentityServer.BaseURL != "rando2" {
		t.Errorf("Wrong identity server received")
	}
}
