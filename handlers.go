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
		fmt.Fprintf(w, authNotification, "is-warning", "already logged in")
		return
	}
	email := r.FormValue("username")
	password := r.FormValue("password")
	u, err := s.GetUserByEmail(email)
	if err != nil {
		s.Logger.Println("user not found", email)
		fmt.Fprintf(w, authNotification, "is-danger", "that straight up did not work")
		return
	}
	ok, err := u.PasswordMatches(password)
	if err != nil {
		s.Logger.Println("error checking password", err)
		fmt.Fprintf(w, authNotification, "is-danger", "that straight up did not work")
		return
	}
	if !ok {
		s.Logger.Println("password does not match", email)
		fmt.Fprintf(w, authNotification, "is-danger", "that straight up did not work")
		return
	}
	u.updateHandle()
	tk, err := u.Token.CreateToken(u.ID, time.Hour)
	if err != nil {
		s.Logger.Println("error creating token", err)
		fmt.Fprintf(w, authNotification, "is-danger", "an error occured when creating token...")
		return
	}
	tk.Email = u.Email
	tk.Handle = u.Handle
	err = s.SaveToken(tk)
	if err != nil {
		s.Logger.Println("error saving token", err)
		fmt.Fprintf(w, authNotification, "is-danger", "an error occured when saving token...")
		return
	}
	err = s.AddTokenToSession(r, w, tk)
	if err != nil {
		s.Logger.Println("error adding token to session", err)
		fmt.Fprintf(w, authNotification, "is-danger", "an error occured when adding token to session...")
		return
	}
	s.Logger.Println("login successful", u.Email)
	fmt.Fprintf(w, authNotification, "is-success", "login successful")
	// theirRoom := fmt.Sprintf("/room/%s", u.ID)
	// w.Header().Set("HX-Redirect", theirRoom)
}

func (s *Server) LoginView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, loginView)
}

func (s *Server) LogoutView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, loginView)
}

func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tk, err := s.GetTokenFromSession(r)
	if err != nil {
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	err = s.DeleteToken(tk)
	if err != nil {
		http.Error(w, "error deleting token", http.StatusInternalServerError)
		return
	}
	err = s.DeleteTokenFromSession(r)
	if err != nil {
		http.Error(w, "error deleting token from session", http.StatusInternalServerError)
		return
	}
	// http.Redirect(w, r, "/access", http.StatusFound)
	w.Header().Set("HX-Redirect", "/access")
}

func (s *Server) clearAuthNotificationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, clearAuthNotification)
}

func (s *Server) MessageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO we could get the user and pass the userid to the WSHMessage
	// gaining the ability to link the div to the user!
	defer func(t time.Time) {
		s.Logger.Println("MessageHandler->time taken: ", time.Since(t))
	}(time.Now())
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
	message = parseCommand(message)
	roomid := r.FormValue("roomid")

	if message == "~clear" {
		room, ok := s.Rooms[roomid]
		if !ok {
			fmt.Println("MessageHandler: room not found", roomid)
			return
		}
		room.ClearMessages()
		message = "hello world!"
	}

	go func(message string, roomid string, token *Token) {
		u, err := s.GetUserByEmail(token.Email)
		if err != nil {
			fmt.Println("MessageHandler: error getting user", err)
			return
		}
		message = SanitizeHTML(message)

		s.Messagechan <- WSMessage{Time: time.Now(), Message: message, Email: token.Handle, RoomID: roomid, UserID: u.ID}
	}(message, roomid, token)
	out := `<input class="input is-outlined" type="text" name="message" id="messageBox" hx-swap-oob="true" placeholder="Type your message...">`
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

	tk, err := s.GetTokenFromSession(r)
	if err != nil {
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	go func(tk string, roomName string) {
		token, err := s.GetToken(tk)
		if err != nil {
			fmt.Println("userHistoryUpdate: error getting token", err)
			return
		}
		u, err := s.GetUserByEmail(token.Email)
		if err != nil {
			fmt.Println("userHistoryUpdate: error getting user", err)
			return
		}
		u.updateHistory(roomName)
		err = s.AddUser(u)
		if err != nil {
			fmt.Println("userHistoryUpdate: error saving user", err)
			return
		}
	}(tk, room.Name)

	fmt.Fprintf(w, chatView, room.ID, room.ID, room.ID, room.Name)

}

func (s *Server) AddRoomToUserRoomsHandler(w http.ResponseWriter, r *http.Request) {
	roomName := r.FormValue("room")
	if roomName == "" {
		http.Error(w, "room name not found", http.StatusBadRequest)
		return
	}
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
	u, err := s.GetUserByEmail(token.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	u.updateRooms(roomName)
	err = s.AddUser(u)
	if err != nil {
		http.Error(w, "error adding user", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "room added")
}

func (s *Server) AddRoomView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, addRoomView)
}

func (s *Server) UserHistoryHandler(w http.ResponseWriter, r *http.Request) {
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
	u, err := s.GetUserByEmail(token.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	out := ""
	encountered := map[string]bool{}
	// tmpl := `<li><a hx-post="/logout" class="has-text-info">logout</a></li>`
	for _, v := range u.History {
		if v == "" {
			continue
		}
		if !encountered[v] {
			encountered[v] = true
			out += fmt.Sprintf(`<li><a href="/room/%s" class="has-text-grey">%s</a></li>`, v, v)
		}

	}
	fmt.Fprint(w, out)

}

func (s *Server) UserRoomsHandler(w http.ResponseWriter, r *http.Request) {
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
	u, err := s.GetUserByEmail(token.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	out := ""
	// tmpl := `<li><a hx-post="/logout" class="has-text-info">logout</a></li>`
	for _, v := range u.Rooms {
		if v == "" {
			continue
		}
		out += fmt.Sprintf(`<li><a href="/room/%s" class="has-text-grey">%s</a></li>`, v, v)
	}
	fmt.Fprint(w, out)
}

func (s *Server) GetRoomStatsHandler(w http.ResponseWriter, r *http.Request) {
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
	u, err := s.GetUserByEmail(token.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	mostRecentRoom := u.History[len(u.History)-1]
	room, err := s.GetRoomByName(mostRecentRoom)
	if err != nil {
		http.Error(w, "error getting room", http.StatusInternalServerError)
		return
	}
	out := room.GetRoomStats()
	fmt.Fprint(w, out)
}

func (s *Server) MessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
	roomName := parts[2]
	room, ok := s.Rooms[roomName]
	if !ok {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}
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

func (s *Server) ProfileView(w http.ResponseWriter, r *http.Request) {
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
	u, err := s.GetUserByEmail(token.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, editProfileView, u.Email, u.FirstName, u.LastName, u.About, u.Email)
}

func (s *Server) AddPostHandler(w http.ResponseWriter, r *http.Request) {
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
	u, err := s.GetUserByEmail(token.Email)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	post := r.FormValue("post")
	u.updatePosts(SanitizeHTML(post))
	err = s.AddUser(u)
	if err != nil {
		http.Error(w, "error saving user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "post added")
}

func (s *Server) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	basicContent := `<div class="content">
		<h1 class="title is-1">Profile</h1>
		<p class="has-text-warning">handle: %v</p>
		<p class="has-text-info">about: %v</p>
		</div>`
	postsContent := `<div class="content">
		<h1 class="title is-1">Posts</h1>
		%v
		</div>`
	postDiv := `<div class="box has-background-black mydisplay">
		<p class="has-text-info">%v</p>
		<a href="/room/%v" target="_blank" rel="noopener noreferrer" class="has-text-white">comment</a>
		</div>`
	tk, err := s.GetTokenFromSession(r)
	if err != nil {
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	_, err = s.GetToken(tk)
	if err != nil {
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	urlParts := strings.Split(r.URL.Path, "/")
	if len(urlParts) < 3 {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	userid := urlParts[2]
	u, err := s.GetUserByID(userid)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}
	posts := ""
	for _, v := range u.Posts {
		posts += fmt.Sprintf(postDiv, v.Content, v.ID)
	}
	posts = fmt.Sprintf(postsContent, posts)
	out := fmt.Sprintf(basicContent, u.Handle, u.About)
	// out += posts
	profileView := fmt.Sprintf(profileView, out, posts)
	fmt.Fprint(w, profileView)

}

func (s *Server) AddPostView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, addPostView)
}

func (s *Server) HelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, "static/scream.pdf")
}

func (s *Server) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.FormValue("userid")
	email := r.FormValue("email")
	// password := r.FormValue("password")
	fname := r.FormValue("first_name")
	lname := r.FormValue("last_name")
	about := r.FormValue("about")
	if len(about) > 200 {
		http.Error(w, "about too long", http.StatusBadRequest)
		return
	}
	u, err := s.GetUserByEmail(userid)
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	u.FirstName = fname
	u.LastName = lname
	u.Email = email
	u.About = about
	u.updateHandle()
	err = s.AddUser(u)
	if err != nil {
		http.Error(w, "error adding user", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "user updated")
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

func parseCommand(c string) string {
	if c == "" {
		return ""
	}
	// ~link__https://foo.bar/thing__link text here
	out := `<a href="%v" class="has-text-link">%v</a>`
	parts := strings.Split(c, "__")
	// fmt.Println(parts, len(parts))
	// if len(parts) < 3 || parts[0] != "~link" {
	// 	return c
	// }
	switch parts[0] {
	case "~link":
		if len(parts) < 3 {
			return c
		}
		return fmt.Sprintf(out, parts[1], parts[2])
	case "~clear":
		return "~clear"
	default:
		return c
	}

}
