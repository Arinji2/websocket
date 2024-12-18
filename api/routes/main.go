package routes

import (
	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	apiRouter := chi.NewRouter()

	// Define POST routes
	apiRouter.Post("/user/create", handleUserCreate)
	apiRouter.Post("/rooms/create", handleRoomCreate)
	apiRouter.Post("/rooms/join", handleRoomJoin)
	apiRouter.Post("/rooms/leave", HandleRoomLeave)

	return apiRouter
}
