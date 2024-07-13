package main

import (
	"net/http"
)

type Room struct {
	ID    string
	Mux   *http.ServeMux
	Reqch chan EnhanchedRequest
	Resch chan *http.Response
}

type EnhanchedRequest struct {
	User *User
	*http.Request
}
