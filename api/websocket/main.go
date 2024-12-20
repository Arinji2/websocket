package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Arinji2/websockets/sqlite"
	"github.com/Arinji2/websockets/types"
	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	Upgrader websocket.Upgrader
	AuthMap  *ClientMap[types.UserData]
	ConnMap  *ClientMap[*websocket.Conn]
	UserMap  *ClientMap[string]
}

func NewWebsocketHandler(authMap *ClientMap[types.UserData]) WebsocketHandler {
	return WebsocketHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return origin == "http://localhost:5173"
			},
		},
		AuthMap: authMap,
		ConnMap: NewClientMap[*websocket.Conn](),
		UserMap: NewClientMap[string](),
	}
}

func (wsh WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}

	connPtr := fmt.Sprintf("%p", c)

	defer func() {
		if userID, exists := wsh.UserMap.Get(connPtr); exists {
			wsh.ConnMap.Delete(userID)
			wsh.UserMap.Delete(connPtr)
			wsh.AuthMap.Delete(userID)
			wsh.Broadcast([]byte(fmt.Sprintf("User %s disconnected", userID)), userID)
		}
		c.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.WriteMessage(websocket.TextMessage, []byte("Connected to server successfully"))
	if err != nil {
		log.Println("Error sending connection message:", err)
		return
	}

	PingHandler(ctx, c)

	authenticated := false
	var userID string

	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				return
			}

			c.SetReadDeadline(time.Now().Add(45 * time.Second))

			if messageType == websocket.BinaryMessage {
				if err := c.WriteMessage(websocket.TextMessage, []byte("Server does not process binary data")); err != nil {
					log.Println("Error sending binary data response:", err)
					return
				}
				continue
			}

			trimmedMessage := strings.TrimSpace(string(message))
			log.Printf("Received message: %s", trimmedMessage)

			if !authenticated {
				if !strings.HasPrefix(trimmedMessage, "AUTH:") {
					c.WriteMessage(websocket.TextMessage, []byte("Please authenticate first"))
					continue
				}

				authMessage := strings.TrimPrefix(trimmedMessage, "AUTH:")
				var authData types.UserData
				if err := json.Unmarshal([]byte(authMessage), &authData); err != nil {
					log.Println("Error unmarshalling auth message:", err)
					return
				}

				if err := sqlite.AuthenticateUser(authData); err != nil {
					log.Println("Error authenticating user:", err)
					c.WriteMessage(websocket.TextMessage, []byte("Authentication failed"))
					return
				}

				userID = authData.ID
				wsh.AuthMap.Add(userID, authData)
				wsh.ConnMap.Add(userID, c)
				wsh.UserMap.Add(connPtr, userID)

				authenticated = true
				log.Printf("User %s authenticated and mapped to connection %s", userID, connPtr)
				wsh.Broadcast([]byte(fmt.Sprintf("New User Connected With Email %s", authData.Email)), userID)
				continue
			}

			switch trimmedMessage {
			case "ping":
				log.Println("Received ping from client")
			case "stop":
				log.Println("Stopping server...")
				cancel()
				return
			default:
				if userData, exists := wsh.AuthMap.Get(userID); exists {
					log.Printf("Message from user: %s", userData.ID)
				}

				err := c.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error sending message: %v", err)
					return
				}
			}
		}
	}
}

// Helper method to broadcast message to all authenticated users except sender
func (wsh WebsocketHandler) Broadcast(message []byte, senderID string) {
	wsh.ConnMap.Range(func(userID string, conn *websocket.Conn) bool {
		if userID != senderID {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error broadcasting to user %s: %v", userID, err)
			}
		}
		return true
	})
}

// Helper method to send message to specific user
func (wsh WebsocketHandler) SendToUser(userID string, message []byte) error {
	if conn, exists := wsh.ConnMap.Get(userID); exists {
		return conn.WriteMessage(websocket.TextMessage, message)
	}
	return fmt.Errorf("user %s not connected", userID)
}
