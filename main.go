package main

import (
	"log"
	"net/http"

	socket "github.com/Arinji2/websockets/websocket"
	"github.com/gorilla/websocket"
)

func main() {
	socketHandler := socket.WebsocketHandler{
		Upgrader: websocket.Upgrader{},
	}
	http.Handle("/", socketHandler)
	log.Println("Server Started on 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
