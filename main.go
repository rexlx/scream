package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	flag.Parse()
	s := NewServer()
	s.Logger.Println("server started")
	fmt.Println("server started")
	if err := http.ListenAndServe(":8080", s.Session.LoadAndSave(s.Gateway)); err != nil {
		s.Logger.Fatalf("error starting server: %v", err)
	}
}
