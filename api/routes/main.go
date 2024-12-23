package routes

import (
	"github.com/Arinji2/websockets/websocket"
	websocketTasks "github.com/Arinji2/websockets/websocket/tasks"
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	apiRouter := chi.NewRouter()
	authMap := websocket.NewClientMap[websocketTasks.UserData]()
	websocketHandler := websocket.NewWebsocketHandler(
		authMap,
	)
	// Define POST routes
	apiRouter.Post("/user/create", handleUserCreate)
	apiRouter.Get("/handler/ws", websocketHandler.ServeHTTP)
	return apiRouter
}
