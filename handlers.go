package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("username")
	password := r.FormValue("password")
	u, err := s.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	ok, err := u.PasswordMatches(password)
	if err != nil {
		http.Error(w, "error comparing passwords", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	tk, err := u.Token.CreateToken(u.ID, time.Hour)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	err = s.SaveToken(tk)
	if err != nil {
		http.Error(w, "error saving token", http.StatusInternalServerError)
		return
	}
	err = s.AddTokenToSession(r, w, tk)
	if err != nil {
		http.Error(w, "error adding token to session", http.StatusInternalServerError)
		return
	}
	successDiv := fmt.Sprintf("<div>%s</div>", "login successful")
	fmt.Fprint(w, successDiv)
}

func (s *Server) LoginView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, loginView)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// /logout
}

func (rm *Room) MessageHandler(w http.ResponseWriter, r *http.Request) {
	// send-message
}

func (rm *Room) PrintMessageHandler(w http.ResponseWriter, r *http.Request) {
	// send-message
}

func (s *Server) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	var u User
	err := u.CreateUser(email, password)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}
	err = s.AddUser(u)
	if err != nil {
		http.Error(w, "error adding user", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "user added")
	w.Header().Set("HX-Redirect", "/access")
	// w.WriteHeader(http.StatusOK)
}

func (s *Server) AddUserView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, AdUserHTML)
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	roomName := getRoomNameFromURL(r.URL.Path)
	if roomName == "" {
		redirectToLogin(w, r)
		return
	}
	s.Memory.RLock()
	room, ok := s.Rooms[roomName]
	s.Memory.RUnlock()
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	room.Mux.ServeHTTP(w, r)

}

func getRoomNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 3 && parts[1] == "room" {
		return parts[2]
	}
	return ""
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}
