package main

import (
	"net/http"
	"time"
)

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.GetTokenFromSession(r)
		if err != nil {
			http.Error(w, "error getting token from session", http.StatusUnauthorized)
			return
		}
		tk, err := s.GetToken(token)
		if err != nil {
			http.Error(w, "error getting token", http.StatusUnauthorized)
			return
		}
		if tk.ExpiresAt.Before(time.Now()) {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
