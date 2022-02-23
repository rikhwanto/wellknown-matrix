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

type headerCombination struct {
	header, value string
}

type urlCombination struct {
	url, value string
}

type serverCombination struct {
	homeserver, identityserver urlCombination
}

var serverList = []serverCombination{
	{homeserver: urlCombination{"CLIENT_HOMESERVER", "rando"},
		identityserver: urlCombination{"CLIENT_IDENTITYSERVER", "rando2"}},
	{homeserver: urlCombination{"CLIENT_HOMESERVER", "rando3"},
		identityserver: urlCombination{"CLIENT_IDENTITYSERVER", "rando4"}},
	{homeserver: urlCombination{"CLIENT_HOMESERVER", "rando5"},
		identityserver: urlCombination{"CLIENT_IDENTITYSERVER", "rando6"}},
	{homeserver: urlCombination{"CLIENT_HOMESERVER", "rando5:444"},
		identityserver: urlCombination{"CLIENT_IDENTITYSERVER", "rando6:4321"}},
}

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

func TestMultipleWithIdentityServer(t *testing.T) {
	for _, testNow := range serverList {
		var headerList = []headerCombination{
			{"Access-Control-Allow-Origin", "*"},
			{"content-type", "application/json;charset=UTF-8"},
		}
		var urlList = []urlCombination{
			testNow.homeserver,
			testNow.identityserver,
		}

		for _, env := range urlList {
			os.Setenv(env.url, env.value)
		}

		req := httptest.NewRequest(http.MethodGet, clientPath, nil)
		w := httptest.NewRecorder()
		requestHandler(w, req)
		result := w.Result()
		if result.StatusCode != http.StatusOK {
			t.Errorf("Wrong status code with status code %d that should be 200", result.StatusCode)
		}

		for _, test := range headerList {
			if headerValue := result.Header.Get(test.header); headerValue != test.value {
				t.Errorf("Wrong header received at client at header %s with response %s that should be %s",
					test.header, headerValue, test.value)
			}
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

		for _, test := range urlList {
			if test.url == "CLIENT_HOMESERVER" {
				if value := clientJSON.Homeserver.BaseURL; value != test.value {
					t.Errorf("Wrong homeserver received, should be %s but got %s", test.value, value)
				}
			} else if test.url == "CLIENT_IDENTITYSERVER" {
				if value := clientJSON.IdentityServer.BaseURL; value != test.value {
					t.Errorf("Wrong identity server received, should be %s but got %s", test.value, value)
				}
			}
		}

	}
}
