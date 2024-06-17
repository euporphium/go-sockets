package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection", err.Error())
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var message Message
		err = conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("Error reading message", err.Error())
			delete(clients, conn)
			return
		}

		broadcast <- message
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				fmt.Println("Error writing message", err.Error())
				client.Close()
				delete(clients, client)
			}
		}
	}
}
