package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

var (
	userBucket  = flag.String("user-bucket", "users", "user bucket")
	tokenBucket = flag.String("token-bucket", "tokens", "token bucket")
	dbName      = flag.String("db-name", "chat.db", "database name")
	logFile     = flag.String("log-file", "chat.log", "log file")
	url         = flag.String("url", ":8080", "url")
	mLimit      = flag.Int("message-limit", 100, "message limit")
)

type Server struct {
	*WSHandler
	Logger  *log.Logger
	DB      *bbolt.DB
	Gateway *http.ServeMux
	Memory  *sync.RWMutex
	Context *context.Context
	Rooms   map[string]*Room
	URL     string
	Session *scs.SessionManager
}

type Token struct {
	Handle    string
	ID        string
	Email     string
	UserID    string
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	Hash      []byte
}

func NewServer() *Server {
	sessionMgr := scs.New()
	sessionMgr.Lifetime = 24 * time.Hour
	sessionMgr.IdleTimeout = 20 * time.Minute
	sessionMgr.Cookie.Persist = true
	sessionMgr.Cookie.Name = "token"
	sessionMgr.Cookie.SameSite = http.SameSiteLaxMode
	// sessionMgr.Cookie.Secure = true
	sessionMgr.Cookie.HttpOnly = true
	fh, err := os.OpenFile(*logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	myLogger := log.New(fh, "", log.LstdFlags)
	mem := &sync.RWMutex{}
	db, err := bbolt.Open(*dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	rooms := make(map[string]*Room)

	wsh := &WSHandler{
		Stop:        make(chan struct{}),
		Conn:        nil,
		Memory:      &sync.RWMutex{},
		Messagechan: make(chan WSMessage, 20),
	}

	s := &Server{
		WSHandler: wsh,
		Logger:    myLogger,
		DB:        db,
		Gateway:   http.NewServeMux(),
		Memory:    mem,
		Context:   nil,
		Rooms:     rooms,
		URL:       *url,
		Session:   sessionMgr,
	}
	// s.Gateway.HandleFunc("/", s.RoomHandler)
	s.Gateway.HandleFunc("/access", s.LoginView)
	s.Gateway.HandleFunc("/login", s.LoginHandler)
	s.Gateway.HandleFunc("/logout", s.LogoutHandler)
	s.Gateway.HandleFunc("/add-user", s.AddUserView)
	s.Gateway.HandleFunc("/adduser", s.AddUserHandler)
	s.Gateway.HandleFunc("/update-profile", s.ProfileHandler)
	s.Gateway.HandleFunc("/can", s.clearAuthNotificationHandler)
	s.Gateway.HandleFunc("/profile", s.ProfileView)
	s.Gateway.Handle("/static/", http.StripPrefix("/static/", s.FileServer()))
	// s.Gateway.HandleFunc("/messagehist", s.MessageHistoryHandler)
	s.Gateway.Handle("/send-message", s.ValidateToken(http.HandlerFunc(s.MessageHandler)))
	s.Gateway.Handle("/ws/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.ServeWS(rooms, w, r)
	}))

	s.Gateway.Handle("/room/", s.ValidateToken(http.HandlerFunc(s.RoomHandler)))
	s.Gateway.Handle("/messagehist/", s.ValidateToken(http.HandlerFunc(s.MessageHistoryHandler)))
	return s
}

func (t *Token) CreateToken(userID string, ttl time.Duration) (*Token, error) {
	tk := &Token{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}
	hotSauce := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, hotSauce)
	if err != nil {
		return nil, err
	}
	tk.Token = uuid.New().String()
	hash := sha256.Sum256([]byte(tk.Token))
	tk.Hash = hash[:]
	return tk, nil
}

func (t *Token) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Token) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (s *Server) TestToken(app *Server, r *http.Request) (bool, error) {
	token, err := s.GetTokenFromSession(r)
	if err != nil {
		return false, err
	}
	tk, err := s.GetToken(token)
	if err != nil {
		return false, err
	}
	if tk.ExpiresAt.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

func (s *Server) GetUserByEmail(email string) (User, error) {
	s.Memory.RLock()
	defer s.Memory.RUnlock()
	var user User
	err := s.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(*userBucket))
		v := b.Get([]byte(email))
		if v == nil {
			s.Logger.Println("user not found")
			return nil
		}

		return user.UnmarshalBinary(v)
	})
	return user, err
}

func (s *Server) AddTokenToSession(r *http.Request, w http.ResponseWriter, tk *Token) error {
	s.Session.Put(r.Context(), "token", tk.Token)
	return nil
}

func (s *Server) DeleteToken(token string) error {
	s.Memory.Lock()
	defer s.Memory.Unlock()
	return s.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(*tokenBucket))
		return b.Delete([]byte(token))
	})
}

func (s *Server) DeleteTokenFromSession(r *http.Request) error {
	s.Session.Remove(r.Context(), "token")
	return nil
}

func (s *Server) FileServer() http.Handler {
	return http.FileServer(http.Dir("./static"))
}

func (s *Server) GetTokenFromSession(r *http.Request) (string, error) {
	tk, ok := s.Session.Get(r.Context(), "token").(string)
	if !ok {
		return "", errors.New("error getting token from session")
	}
	return tk, nil
}

func (s *Server) AddUser(u User) error {
	s.Memory.Lock()
	defer s.Memory.Unlock()
	return s.DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(*userBucket))
		if err != nil {
			return err
		}
		v, err := u.MarshalBinary()
		if err != nil {
			return err
		}
		return b.Put([]byte(u.Email), v)
	})
}

func (s *Server) SaveToken(tk *Token) error {
	s.Memory.Lock()
	defer s.Memory.Unlock()
	return s.DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(*tokenBucket))
		if err != nil {
			return err
		}
		v, err := tk.MarshalBinary()
		if err != nil {
			return err
		}
		return b.Put([]byte(tk.Token), v)
	})
}

func (s *Server) GetToken(token string) (*Token, error) {
	s.Memory.RLock()
	defer s.Memory.RUnlock()
	var tk Token
	err := s.DB.View(func(tx *bbolt.Tx) error {
		// b, err := tx.CreateBucketIfNotExists([]byte(*tokenBucket))
		// if err != nil {
		// 	return err
		// }
		b := tx.Bucket([]byte(*tokenBucket))
		v := b.Get([]byte(token))
		if v == nil {
			return nil
		}
		return tk.UnmarshalBinary(v)
	})
	return &tk, err
}

func (s *Server) AddRoom(r *Room) {
	s.Memory.Lock()
	defer s.Memory.Unlock()
	s.Rooms[r.ID] = r
}
