package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server struct {
	StartTime time.Time
	Logger    *log.Logger
	Memory    *sync.RWMutex
	Gateway   *http.ServeMux
}

func NewServer(fh string) (*Server, error) {
	file, err := os.OpenFile(fh, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	l := log.New(file, "", log.LstdFlags)
	s := &Server{
		StartTime: time.Now(),
		Logger:    l,
		Memory:    &sync.RWMutex{},
		Gateway:   http.NewServeMux(),
	}
	s.Gateway.HandleFunc("/graph", s.CreateGraph)

	return s, nil
}
