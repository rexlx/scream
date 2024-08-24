package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"go.etcd.io/bbolt" // this shows an error even though it is installed
)

var (
	userBucket    = flag.String("user-bucket", "users", "user bucket")
	tokenBucket   = flag.String("token-bucket", "tokens", "token bucket")
	dbName        = flag.String("db-name", "chat.db", "database name")
	logFile       = flag.String("log-file", "chat.log", "log file")
	url           = flag.String("url", ":8081", "url")
	mLimit        = flag.Int("message-limit", 100, "message limit")
	certFile      = flag.String("cert-file", "server-cert.pem", "cert file")
	keyFile       = flag.String("key-file", "server-key.pem", "key file")
	firstUserMode = flag.Bool("first-user-mode", false, "first user mode")
	updateFreq    = flag.Duration("update-freq", 2*time.Minute, "update frequency")
)

type Server struct {
	*WSHandler
	Logger      *log.Logger
	DB          *bbolt.DB
	Gateway     *http.ServeMux
	Memory      *sync.RWMutex
	Context     *context.Context
	Rooms       map[string]*Room
	Stats       map[string]float64
	Coordinates map[string][]float64
	Graphs      map[string]string
	GraphCache  string
	URL         string
	Session     *scs.SessionManager
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

// TODO were not actually storing stats in any sort of time series based way
// i think literally just a single point in time lol
func NewServer(url string, firstUser bool) *Server {
	defer func(t time.Time) {
		fmt.Println("NewServer: time taken: ", time.Since(t))
	}(time.Now())
	sessionMgr := scs.New()
	sessionMgr.Lifetime = 24 * time.Hour
	sessionMgr.IdleTimeout = 1 * time.Hour
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

	stats := make(map[string]float64)
	graphs := make(map[string]string)
	coords := make(map[string][]float64)
	wsh := &WSHandler{
		TTL:         2 * time.Hour,
		Stop:        make(chan struct{}),
		Conn:        nil,
		Memory:      &sync.RWMutex{},
		Messagechan: make(chan WSMessage, 20),
	}

	s := &Server{
		Graphs:      graphs,
		Coordinates: coords,
		Stats:       stats,
		WSHandler:   wsh,
		Logger:      myLogger,
		DB:          db,
		Gateway:     http.NewServeMux(),
		Memory:      mem,
		Context:     nil,
		Rooms:       rooms,
		URL:         url,
		Session:     sessionMgr,
	}
	// s.Gateway.HandleFunc("/", s.RoomHandler)
	s.Gateway.HandleFunc("/access", s.LoginView)
	s.Gateway.HandleFunc("/login", s.LoginHandler)
	s.Gateway.HandleFunc("/logout", s.LogoutHandler)
	s.Gateway.HandleFunc("/add-room", s.AddRoomView)
	s.Gateway.HandleFunc("/addroom", s.AddRoomToUserRoomsHandler)
	s.Gateway.HandleFunc("/adduser", s.AddUserHandler)
	s.Gateway.HandleFunc("/addpost", s.AddPostHandler)
	s.Gateway.HandleFunc("/add-post", s.AddPostView)
	s.Gateway.HandleFunc("/update-profile", s.ProfileHandler)
	s.Gateway.HandleFunc("/can", s.clearAuthNotificationHandler)
	s.Gateway.HandleFunc("/profile", s.ProfileView)
	s.Gateway.HandleFunc("/history", s.UserHistoryHandler)
	s.Gateway.HandleFunc("/rooms", s.UserRoomsHandler)
	s.Gateway.HandleFunc("/help", s.HelpHandler)
	s.Gateway.HandleFunc("/roomstats", s.GetRoomStatsHandler)
	s.Gateway.Handle("/static/", http.StripPrefix("/static/", s.FileServer()))
	// s.Gateway.HandleFunc("/messagehist", s.MessageHistoryHandler)
	s.Gateway.Handle("/send-message", s.ValidateToken(http.HandlerFunc(s.MessageHandler)))
	s.Gateway.Handle("/stats", s.ValidateToken(http.HandlerFunc(s.StatHandler)))
	if firstUser {
		s.Gateway.HandleFunc("/add-user", s.AddUserView)
	} else {
		s.Gateway.Handle("/add-user", s.ValidateToken(http.HandlerFunc(s.AddUserView)))
	}
	s.Gateway.Handle("/splash", s.ValidateToken(http.HandlerFunc(s.SplashView)))
	s.Gateway.Handle("/ws/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.ServeWS(rooms, w, r)
	}))

	s.Gateway.Handle("/room/", s.ValidateToken(http.HandlerFunc(s.RoomHandler)))
	s.Gateway.Handle("/user/", s.ValidateToken(http.HandlerFunc(s.UserProfileHandler)))
	s.Gateway.Handle("/messagehist/", s.ValidateToken(http.HandlerFunc(s.MessageHistoryHandler)))
	s.Gateway.Handle("/", http.HandlerFunc(s.LoginView))
	// s.Stats["start"] = float64(time.Now().Unix())
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
	s.Memory.Lock()
	defer s.Memory.Unlock()
	s.Stats["user_queries"]++
	var user User
	err := s.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(*userBucket))
		v := b.Get([]byte(email))
		if v == nil {
			s.Logger.Println("user not found")
			s.Stats["user_not_found"]++
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
	u.UpdatedAt = time.Now()
	s.Memory.Lock()
	defer s.Memory.Unlock()
	s.Stats["user_saves"]++
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
	s.Stats["token_saves"]++
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
		b := tx.Bucket([]byte(*tokenBucket))
		v := b.Get([]byte(token))
		if v == nil {
			return nil
		}
		s.Stats["token_queries"]++
		return tk.UnmarshalBinary(v)
	})
	return &tk, err
}

func (s *Server) AddRoom(r *Room) {
	s.Memory.Lock()
	defer s.Memory.Unlock()
	s.Stats["room_created"]++
	s.Rooms[r.ID] = r
}

func SanitizeHTML(s string) string {
	s = html.EscapeString(s)
	return s
}

func (s *Server) GetUserByID(userid string) (User, error) {
	s.Memory.RLock()
	defer s.Memory.RUnlock()
	var user User
	err := s.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(*userBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := user.UnmarshalBinary(v)
			if err != nil {
				return err
			}
			if user.ID == userid {
				return nil
			}
		}
		return nil
	})
	return user, err
}

func (s *Server) CleanUpTokens() error {
	s.Memory.Lock()
	defer s.Memory.Unlock()
	return s.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(*tokenBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tk Token
			err := tk.UnmarshalBinary(v)
			if err != nil {
				return err
			}
			if tk.ExpiresAt.Before(time.Now()) {
				err := b.Delete(k)
				if err != nil {
					return err
				}
				fmt.Println("token deleted", tk.Token)
			}
		}
		return nil
	})
}

func (s *Server) UpdateGraphs(t time.Duration) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	s.Memory.Lock()
	malloc := float64(m.Alloc)
	s.Stats["malloc"] = malloc / 1024
	s.Stats["goroutines"] = float64(runtime.NumGoroutine())
	s.Stats["heap"] = float64(m.HeapAlloc) / 1024
	s.Stats["heap_objects"] = float64(m.HeapObjects)
	s.Stats["stack"] = float64(m.StackInuse) / 1024
	s.Stats["alloc"] = float64(m.Alloc) / 1024
	s.Stats["total_alloc"] = float64(m.TotalAlloc) / 1024
	s.Stats["sys"] = float64(m.Sys) / 1024
	s.Stats["num_gc"] = float64(m.NumGC)
	s.Stats["poll_time"] = float64(time.Now().Unix())
	s.Stats["poll_interval"] = float64(t.Seconds())
	// s.Stats["last_gc"] = float64(m.LastGC) / 1000000
	s.Stats["pause_total_ns"] = float64(m.PauseTotalNs) / 1000000
	for i, e := range s.Stats {
		_, ok := s.Coordinates[i]
		if !ok {
			s.Coordinates[i] = make([]float64, 0)
		}
		if len(s.Coordinates[i]) > 99 {
			s.Coordinates[i] = s.Coordinates[i][1:]
		}
		s.Coordinates[i] = append(s.Coordinates[i], e)
		// s.Graphs[i] = s.CreatePolyline(i, s.Coordinates[i], 180, 60)
	}
	s.Memory.Unlock()
	client := http.Client{}
	out, err := json.Marshal(s.Coordinates)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:8081/graph", bytes.NewBuffer(out))
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var resData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resData)
	if err != nil {
		fmt.Println(err)
		return
	}
	val, ok := resData["chart"]
	if !ok {
		fmt.Println("no chart in response")
		return
	}
	s.GraphCache = val.(string)
}

func (s *Server) CreatePolyline(name string, data []float64, width, height int) string {
	var coords string
	maxVlaue := findMax(data)
	scaleX := float64(width) / float64(len(data))
	scaleY := float64(height) / maxVlaue
	for i, v := range data {
		x := float64(i) * scaleX
		y := float64(height) - (v * scaleY)
		coords += fmt.Sprintf("%f,%f ", x, y)
		fmt.Printf("i %v v %v -> x %v y %v\n", i, v, x, y)
	}
	return fmt.Sprintf(polyLineSVG, width, height, coords)
}

func findMax(data []float64) float64 {
	var max float64
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}
