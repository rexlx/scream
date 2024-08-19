package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSMessage struct {
	RoomID  string    `json:"room_id"`
	Time    time.Time `json:"time"`
	From    string    `json:"from"`
	Message string    `json:"message"`
	UserID  string    `json:"user_id"`
	Email   string    `json:"email"`
}

type WSHandler struct {
	TTL         time.Duration
	Stop        chan struct{}
	Conn        *websocket.Conn
	Memory      *sync.RWMutex
	Messagechan chan WSMessage
}

func (wsh *WSHandler) Write(rooms map[string]*Room) {
	var lastMessage time.Time
	var ticker = time.NewTicker(wsh.TTL)
	fmt.Println("WSHandler.Write: new writer")
	defer wsh.Conn.Close()
	defer fmt.Println("WSHandler.Write: closing connection")
dasWriter:
	for {
		select {
		case <-ticker.C:
			if time.Since(lastMessage) > wsh.TTL {
				fmt.Println("WSHandler.Write: closing connection due to inactivity")
				break dasWriter
			}
		case message := <-wsh.Messagechan:
			lastMessage = time.Now()
			// fmt.Printf("got message %+v", message)
			room, ok := rooms[message.RoomID]
			if !ok {
				fmt.Println("WSHandler.Write: room not found", message.RoomID)
				continue
			}
			room.AddMessage(message)
			room.Memory.RLock()
			out := room.GetMesssages()
			for conn := range room.Connections {
				err := conn.WriteMessage(websocket.TextMessage, []byte(out))
				if err != nil {
					fmt.Println("WSHandler.Write: error writing message", err)
					conn.Close()
					delete(room.Connections, conn)
				}
			}
			room.Memory.RUnlock()
		case <-wsh.Stop:
			break dasWriter

		}
	}
}

func (wsh *WSHandler) ServeWS(rooms map[string]*Room, w http.ResponseWriter, r *http.Request) {
	parts := r.URL.Path
	roomID := parts[len("/ws/"):]
	if roomID == "" {
		http.Error(w, "room id not found", http.StatusBadRequest)
		return
	}

	room, ok := rooms[roomID]
	if !ok {
		http.Error(w, "room id not found", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "error upgrading connection", http.StatusInternalServerError)
		return
	}
	wsh.Conn = conn

	room.AddConnection(wsh)

	go wsh.Write(rooms)

}
