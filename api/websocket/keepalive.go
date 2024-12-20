package websocket

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func PingHandler(ctx context.Context, c *websocket.Conn) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					log.Printf("Error sending ping: %v", err)
					return
				}
				log.Println("Sent ping")
			}
		}
	}()

	c.SetReadDeadline(time.Now().Add(45 * time.Second))
}
