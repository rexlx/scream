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
	Stop        chan struct{}
	Conn        *websocket.Conn
	Memory      *sync.RWMutex
	Messagechan chan WSMessage
}

func (wsh *WSHandler) Read() {
	defer wsh.Conn.Close()
	for {
		_, message, err := wsh.Conn.ReadMessage()
		if err != nil {
			break
		}
		m := WSMessage{Time: time.Now(), Message: string(message)}
		// fmt.Println(m)
		wsh.Messagechan <- m
	}
}

func (wsh *WSHandler) Write(rooms map[string]*Room) {
	defer wsh.Conn.Close()
dasWriter:
	for {
		select {
		case message := <-wsh.Messagechan:
			room, ok := rooms[message.RoomID]
			if !ok {
				fmt.Println("room not found", message.RoomID)
				continue
			}
			room.AddMessage(message)
			out := room.GetMesssages()
			room.Memory.RLock()
			// fmt.Println("writing message", out)
			for _, conn := range room.Connections {
				// if conn.Conn == wsh.Conn {
				// 	fmt.Println("wont write to self")
				// 	continue
				// }
				// fmt.Println("writing to conn", conn.Conn)
				err := conn.Conn.WriteMessage(websocket.TextMessage, []byte(out))
				if err != nil {
					fmt.Println("error writing message", err)
					continue
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
	room.AddConnection(wsh)

	fmt.Println("serving ws")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "error upgrading connection", http.StatusInternalServerError)
		return
	}
	wsh.Conn = conn

	go wsh.Write(rooms)

}
