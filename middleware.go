package main

import (
	"net/http"
	"time"
)

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.GetTokenFromSession(r)
		if err != nil {
			s.Logger.Println("error getting token from session", err)
			w.Header().Set("HX-Redirect", "/access")
			http.Error(w, "error getting token from session", http.StatusUnauthorized)
			return
		}
		tk, err := s.GetToken(token)
		if err != nil {
			s.Logger.Println("error getting token", token)
			http.Error(w, "error getting token", http.StatusUnauthorized)
			return
		}
		if tk.ExpiresAt.Before(time.Now()) {
			s.Logger.Println("token expired", tk.ExpiresAt)
			w.Header().Set("HX-Redirect", "/access")
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}
		s.Logger.Println("token valid")
		next.ServeHTTP(w, r)
	})
}
