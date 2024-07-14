package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	_tk, _ := s.GetTokenFromSession(r)
	if _tk != "" {
		http.Error(w, "already logged in", http.StatusUnauthorized)
		return
	}
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
	tk.Email = u.Email
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
	message := r.FormValue("message")
	fmt.Fprintf(w, "message: %s", message)
}

func (rm *Room) PrintMessageHandler(w http.ResponseWriter, r *http.Request) {
	// send-message
}

func (rm *Room) ChatView(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, chatView)
	fmt.Fprintf(w, chatView, rm.ID)
}

func (s *Server) MessageHandler(w http.ResponseWriter, r *http.Request) {
	tk, err := s.GetTokenFromSession(r)
	if err != nil {
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	token, err := s.GetToken(tk)
	if err != nil {
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	message := r.FormValue("message")
	roomid := r.FormValue("roomid")
	out := `<div class="control is-expanded">
          <input class="input is-outlined" type="text" name="message" placeholder="Type your message...">
        </div>
        <div class="control">
          <button class="button is-info is-outlined" type="submit">Send</button>
		  <small class="has-text-grey-light">message sent</small>
        </div>`
	s.Messagechan <- WSMessage{Time: time.Now(), Message: message, Email: token.Email, RoomID: roomid}
	fmt.Fprint(w, out)
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

func (s *Server) RoomHandler(w http.ResponseWriter, r *http.Request) {
	roomName := getRoomNameFromURL(r.URL.Path)
	if roomName == "" {
		redirectToLogin(w, r)
		return
	}
	room, err := s.GetRoomByName(roomName)
	if err != nil {
		room = NewRoom(roomName, *mLimit)
		s.AddRoom(room)
	}

	fmt.Fprintf(w, chatView, room.ID, room.ID, room.ID)

}

func (s *Server) MessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	fmt.Println(parts)
	if len(parts) < 3 {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	roomName := parts[2]
	fmt.Println("roomName", roomName)
	room, ok := s.Rooms[roomName]
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	fmt.Println(room.GetMesssages())
	fmt.Fprint(w, room.GetMesssages())
}

func (s *Server) GetRoomByName(name string) (*Room, error) {
	s.Memory.RLock()
	defer s.Memory.RUnlock()
	for k, v := range s.Rooms {
		if v.Name == name {
			return s.Rooms[k], nil
		}

	}
	return nil, fmt.Errorf("GetRoomByName: room not found")
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
