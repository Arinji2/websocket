package socket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	Upgrader websocket.Upgrader
}

func (wsh WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error In Upgrading Connection To Websocket", err)
		return

	}
	cx, cancel := context.WithCancel(r.Context())

	defer func() {
		fmt.Println("Closing Connection")
		c.Close()
		cancel()
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte("connected to server successfully"))
	if err != nil {
		log.Printf("Error sending message %v", err)
		return
	}
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message %v", err)
			return
		}
		if messageType == websocket.BinaryMessage {
			err = c.WriteMessage(websocket.TextMessage, []byte("server does not process binary data"))
			if err != nil {
				log.Printf("Error sending message %v", err)
			}
		}

		trimmedMessage := strings.Trim(string(message), "\n")
		if trimmedMessage == "ping" {
			fmt.Println("Server Pinged")
			c.WriteMessage(websocket.TextMessage, []byte("SERVER: pong"))
		}

		if trimmedMessage == "start" {
			go func(ctx context.Context) {
				i := 1
				c.WriteMessage(websocket.TextMessage, []byte("SERVER: Starting Continuous Pinging"))
				for {
					select {
					case <-cx.Done():
					default:
						err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("SERVER: ping x%d", i)))
						if err != nil {
							log.Printf("Error sending message %v", err)
							return
						}
					}
					time.Sleep(2 * time.Second)
					i++
				}
			}(cx)
		}

		if trimmedMessage == "stop" {
			fmt.Println("Stopping Pinging")
			cancel()
			c.WriteMessage(websocket.TextMessage, []byte("SERVER: Stopping Pinging"))
		}
	}
}
