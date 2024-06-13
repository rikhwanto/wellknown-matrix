package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type jsonURL struct {
	BaseURL string `json:"base_url,omitempty"`
	RawURL  string `json:"url,omitempty"`
}

type ClientResponse struct {
	Homeserver     jsonURL  `json:"m.homeserver"`
	IdentityServer *jsonURL `json:"m.identity_server,omitempty"`
}

type ServerResponse struct {
	Server string `json:"m.server"`
}

func composeServerResponse() string {
	server := ServerResponse{
		Server: os.Getenv("FEDERATION_SERVER"),
	}
	responseJson, err := json.Marshal(server)
	if err != nil {
		log.Fatal(err)
	}

	return string(responseJson)
}

func composeClientResponse() string {
	homeserverURL := os.Getenv("CLIENT_HOMESERVER")
	identityServerURL := os.Getenv("CLIENT_IDENTITYSERVER")
	client := ClientResponse{
		Homeserver: jsonURL{BaseURL: homeserverURL},
	}
	if identityServerURL != "" {
		client.IdentityServer = &jsonURL{BaseURL: identityServerURL}
	}
	responseJson, err := json.Marshal(client)
	if err != nil {
		log.Fatal(err)
	}
	return string(responseJson)
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	clientPath := "/.well-known/matrix/client"
	serverPath := "/.well-known/matrix/server"
	if req.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if req.URL.Path == serverPath {
		w.Header().Set("content-type", "application/json;charset=UTF-8")
		fmt.Fprint(w, composeServerResponse())
		log.Printf("Served a server well-known page at %s\n", req.URL.Path)
	} else if req.URL.Path == clientPath {
		w.Header().Set("content-type", "application/json;charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, composeClientResponse())
		log.Printf("Served a client well-known page at %s\n", req.URL.Path)
	}
}

func main() {
	http.HandleFunc("/.well-known/matrix/", requestHandler)
	log.Println("Starting a server to serve .well-known files at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
