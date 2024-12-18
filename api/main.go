package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Arinji2/websockets/routes"
	"github.com/Arinji2/websockets/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Server Started on Port 8080")
	db, err := sqlite.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB Found and Ready")
	db.Close()
	r := chi.NewRouter()
	frontendURL := os.Getenv("FRONTEND_URL")
	fmt.Println("Allowing Origins: ", frontendURL)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Mount("/api", routes.Router())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Landing Backend: Health Check"))
		render.Status(r, 200)
	})

	if err := http.ListenAndServe(":8080", r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
