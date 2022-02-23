package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
)

type headerCombination struct {
	header, value string
}

type urlCombination struct {
	url, value string
}

type serverCombination struct {
	homeserver, identityserver urlCombination
}

func generateWrongPaths(limit int) []string {
	var result = []string{"/"}
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= limit; i++ {
		perLimit := rand.Intn(25)
		for j := 0; j < perLimit; j++ {
			result = append(result, fmt.Sprintf("/%s", petname.Generate(i, "/")))
		}

	}
	return result
}

func generateServerLists(limit int) []serverCombination {
	var result []serverCombination
	rand.Seed(time.Now().UnixNano())
	var protocols = []string{"https", "http"}
	envHome := "CLIENT_HOMESERVER"
	envIdentity := "CLIENT_IDENTITYSERVER"
	for i := 0; i < limit; i++ {
		for _, homeProtocol := range protocols {
			for _, identityProtocol := range protocols {
				homeDomain := petname.Generate(2, ".")
				identityDomain := petname.Generate(2, ".")
				temp := serverCombination{
					homeserver:     urlCombination{envHome, fmt.Sprintf("%s://%s", homeProtocol, homeDomain)},
					identityserver: urlCombination{envIdentity, fmt.Sprintf("%s://%s", identityProtocol, identityDomain)},
				}
				result = append(result, temp)
			}
		}
	}
	return result
}

var clientPath = "/.well-known/matrix/client"
var serverPath = "/.well-known/matrix/server"

func TestWrongPath(t *testing.T) {
	wrongPaths := generateWrongPaths(10)
	for _, test := range wrongPaths {
		req := httptest.NewRequest(http.MethodGet, test, nil)
		w := httptest.NewRecorder()
		requestHandler(w, req)
		result := w.Result()
		t.Logf("Received status code %d at wrong path %s", result.StatusCode, test)
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("Wrong status code at wrong path at path %s with status code %d", test, result.StatusCode)
		}
	}
}

func TestWrongMethod(t *testing.T) {
	var wrongMethods = []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
	for _, test := range wrongMethods {
		req := httptest.NewRequest(test, clientPath, nil)
		w := httptest.NewRecorder()
		requestHandler(w, req)
		result := w.Result()
		t.Logf("Received status code %d at wrong method %s", result.StatusCode, test)
		if result.StatusCode != http.StatusNotFound {
			t.Errorf("Wrong status code at wrong method at method %s with status code %d", test, result.StatusCode)
		}
	}
}

func TestClient(t *testing.T) {
	clientServerList := generateServerLists(15)
	for _, testNow := range clientServerList {
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
