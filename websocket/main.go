package socket

import (
	"fmt"
	"net/http"

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
	defer func() {
		fmt.Println("Closing Connection")
		c.Close()
	}()
}
