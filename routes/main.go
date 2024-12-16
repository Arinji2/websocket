package routes

import (
	"log"
	"net/http"

	"github.com/Arinji2/websockets/sqlite"
)

func Router() *http.ServeMux {
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("POST /user/create", handleUserCreate)
	return apiRouter
}

func handleUserCreate(w http.ResponseWriter, r *http.Request) {
	/*
		curl -X POST http://localhost:8080/api/user/create -H "Content-Type: application/json" -d '{"name": "arinji", "email": "arinjaydhar205@gmail.com"}'
	*/
	var userData UserDataRoute
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
	userID, err := GenerateID(r.Context(), db, "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error generating ID: %v", err)
		return
	}

	user := UserData{
		ID:    userID,
		Email: userData.Email,
		Name:  userData.Name,
	}

	_, err = db.Exec(r.Context(), "INSERT INTO users (id, email, name) VALUES (?, ?, ?)", user.ID, user.Email, user.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting user: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user.ID))
	log.Printf("User created with ID: %s For Email: %s", user.ID, user.Email)
}
