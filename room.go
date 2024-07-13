package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type Room struct {
	Connections []*WSHandler
	Memory      *sync.RWMutex
	Messages    []WSMessage
	ID          string
	Gateway     *http.ServeMux
}

type EnhanchedRequest struct {
	User *User
	*http.Request
}

func NewRoom() *Room {
	uid := uuid.New().String()
	mem := &sync.RWMutex{}
	r := &Room{
		ID:      uid,
		Gateway: http.NewServeMux(),
		Memory:  mem,
	}
	r.Gateway.HandleFunc("/send-message", r.MessageHandler)
	r.Gateway.HandleFunc("/messages", r.PrintMessageHandler)
	// r.Gateway.Handle("/room/", http.StripPrefix("/room", http.HandlerFunc(r.ChatView)))
	r.Gateway.Handle("/", http.HandlerFunc(r.ChatView))
	return r

}

func (rm *Room) AddMessage(msg WSMessage) {
	rm.Memory.Lock()
	defer rm.Memory.Unlock()
	rm.Messages = append(rm.Messages, msg)
}

func (rm *Room) GetMesssages() string {
	out := `<div class="box has-background-dark" id="chat-box">
	  %v
        </div>`
	rm.Memory.RLock()
	defer rm.Memory.RUnlock()
	messages := ""
	for _, msg := range rm.Messages {
		messages += fmt.Sprintf(`<em class="mb2 has-text-link">%v</em>:  <p class="mb2">%v</p>`, msg.Email, msg.Message)
	}
	out = fmt.Sprintf(out, messages)
	return out
}

func (rm *Room) AddConnection(conn *WSHandler) {
	rm.Memory.Lock()
	defer rm.Memory.Unlock()
	fmt.Println("adding connection", conn)
	rm.Connections = append(rm.Connections, conn)
}

func (rm *Room) MessagesHandler(w http.ResponseWriter, r *http.Request) {

}
