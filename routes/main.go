package routes

import (
	"net/http"
)

func Router() *http.ServeMux {
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("POST /user/create", handleUserCreate)
	return apiRouter
}
