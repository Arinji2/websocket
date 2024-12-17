package routes

import (
	"net/http"
)

func Router() *http.ServeMux {
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("POST /user/create", handleUserCreate)
	apiRouter.HandleFunc("POST /rooms/create", handleRoomCreate)
	apiRouter.HandleFunc("POST /rooms/join", handleRoomJoin)
	apiRouter.HandleFunc("POST /rooms/leave", HandleRoomLeave)
	return apiRouter
}
