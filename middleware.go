package main

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.GetTokenFromSession(r)
		if err != nil {
			fmt.Println("error getting token from session")
			http.Error(w, "error getting token from session", http.StatusUnauthorized)
			return
		}
		tk, err := s.GetToken(token)
		if err != nil {
			fmt.Println("error getting token", token)
			http.Error(w, "error getting token", http.StatusUnauthorized)
			return
		}
		if tk.ExpiresAt.Before(time.Now()) {
			fmt.Println("token expired")
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}
		fmt.Println("token valid")
		next.ServeHTTP(w, r)
	})
}
