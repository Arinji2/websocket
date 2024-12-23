package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Arinji2/websockets/sqlite"
	websocketTasks "github.com/Arinji2/websockets/websocket/tasks"
)

func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	/*
		curl -X POST http://localhost:8080/api/user/create -H "Content-Type: application/json" -d '{"name": "arinji", "email": "arinjaydhar205@gmail.com"}'
	*/
	var userData websocketTasks.UserDataRoute
	if err := parseRequestBody(r, &userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sqlite.NewConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating database connection: %v", err)
		return
	}

	defer db.Close()
	userID, err := sqlite.GenerateID(r.Context(), db, "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error generating ID: %v", err)
		return
	}

	sessionID := generateSessionID()
	user := websocketTasks.UserData{
		ID:        userID,
		Email:     userData.Email,
		Name:      userData.Name,
		SessionID: sessionID,
	}

	_, err = db.Exec(r.Context(), "INSERT INTO users (id, email, name, session_id) VALUES (?, ?, ?, ?)", user.ID, user.Email, user.Name, user.SessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting user: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	log.Printf("User created with ID: %s For Email: %s", user.ID, user.Email)
}
