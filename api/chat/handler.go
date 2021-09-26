package chat

import (
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Chat struct {
	DB  *gorm.DB
	Hub *Hub
}

func Routes(chat *Chat) *Chat {
	return chat
}

func (chat *Chat) Socket() *Hub {
	hub := Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	chat.Hub = &hub
	return chat.Hub
}

func (chat *Chat) HandleChat(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: chat.Hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
