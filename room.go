package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Name         string
	MessageLimit int
	ID           string
	Messages     []WSMessage
	Connections  map[*websocket.Conn]bool
	Memory       *sync.RWMutex
	// Gateway      *http.ServeMux
}

type EnhanchedRequest struct {
	User *User
	*http.Request
}

func NewRoom(name string, mLimit int) *Room {
	uid := uuid.New().String()
	mem := &sync.RWMutex{}
	conns := make(map[*websocket.Conn]bool)
	r := &Room{
		Connections:  conns,
		Name:         name,
		MessageLimit: mLimit,
		ID:           uid,
		// Gateway:      http.NewServeMux(),
		Memory: mem,
	}
	// r.Gateway.HandleFunc("/send-message", r.MessageHandler)
	// r.Gateway.HandleFunc("/messages", r.PrintMessageHandler)
	// r.Gateway.Handle("/room/", http.StripPrefix("/room", http.HandlerFunc(r.ChatView)))
	// r.Gateway.Handle("/", http.HandlerFunc(r.ChatView))
	return r

}

func (rm *Room) GetRoomStats() string {
	var div = `<div class="content" id="roomstats">
	%v
</div>`
	var c int
	rm.Memory.RLock()
	defer rm.Memory.RUnlock()
	for conn := range rm.Connections {
		if conn != nil {
			c++
		}
	}
	out := `<strong>%v</strong>: %v users, %v messages`
	out = fmt.Sprintf(out, rm.Name, c, len(rm.Messages))
	div = fmt.Sprintf(div, out)
	return div
}

func (rm *Room) AddMessage(msg WSMessage) {
	rm.Memory.Lock()
	defer rm.Memory.Unlock()
	if len(rm.Messages) >= rm.MessageLimit {
		rm.Messages = rm.Messages[1:]
	}
	rm.Messages = append(rm.Messages, msg)
}

func (rm *Room) GetMesssages() string {
	out := `<div class="box has-background-black mydisplay" id="chat-box">
	  %v
        </div>`
	rm.Memory.RLock()
	defer rm.Memory.RUnlock()
	messages := ""
	for _, msg := range rm.Messages {
		messages += fmt.Sprintf(`<div class="content has-background-black"><em class="has-text-white">%v:</em>  <p class="has-text-primary">%v</p></div>`, msg.Email, msg.Message)
	}
	if messages == "" {
		messages = `<div class="content has-background-black"><em class="has-text-white">server bot:</em>  <p class="has-text-primary">you're the first one here! (maybe)</p></div>`
	}
	out = fmt.Sprintf(out, messages)
	return out
}

func (rm *Room) AddConnection(conn *WSHandler) {
	rm.Memory.Lock()
	defer rm.Memory.Unlock()
	rm.Connections[conn.Conn] = true
}

func (rm *Room) MessagesHandler(w http.ResponseWriter, r *http.Request) {

}
