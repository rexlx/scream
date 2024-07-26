package main

import (
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Rooms     []string  `json:"rooms"`
	Posts     []string  `json:"posts"`
	History   []string  `json:"history"`
	About     string    `json:"about"`
	ID        string    `json:"id"`
	Handle    string    `json:"handle"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     Token     `json:"token"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *User) updateHandle() {
	u.Handle = u.FirstName + u.LastName
	if u.Handle == "" {
		u.Handle = u.Email
	}
}

func (u *User) updateHistory(roomid string) {
	if len(u.History) >= 10 {
		u.History = u.History[1:]
	}
	u.History = append(u.History, roomid)
}

func (u *User) updatePosts(postid string) {
	if len(u.Posts) >= 10 {
		u.Posts = u.Posts[1:]
	}
	u.Posts = append(u.Posts, postid)
}

func (u *User) updateRooms(roomid string) {
	if len(u.Rooms) >= 10 {
		u.Rooms = u.Rooms[1:]
	}
	u.Rooms = append(u.Rooms, roomid)
}

func (u *User) PasswordMatches(input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input))
	if err != nil {
		return false, err
	}
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			//invalid password
			return false, nil
		default:
			//unknown error
			return false, err
		}
	}
	return true, nil
}

func (u *User) CreateUser(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Email = email
	u.Password = string(hash)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.History = make([]string, 0)
	u.Rooms = make([]string, 0)
	u.Posts = make([]string, 0)
	return nil
}
