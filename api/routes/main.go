package routes

import (
	"github.com/Arinji2/websockets/types"
	"github.com/Arinji2/websockets/websocket"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	apiRouter := chi.NewRouter()
	authMap := websocket.NewClientMap[types.UserData]()
	websocketHandler := websocket.NewWebsocketHandler(
		authMap,
	)
	// Define POST routes
	apiRouter.Post("/user/create", handleUserCreate)
	apiRouter.Post("/rooms/create", handleRoomCreate)
	apiRouter.Post("/rooms/join", handleRoomJoin)
	apiRouter.Post("/rooms/leave", HandleRoomLeave)
	apiRouter.Get("/handler/ping", websocketHandler.ServeHTTP)
	return apiRouter
}
