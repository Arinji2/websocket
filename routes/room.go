package routes

import (
	"log"
	"net/http"

	"github.com/Arinji2/websockets/sqlite"
)

func handleRoomCreate(w http.ResponseWriter, r *http.Request) {
	/*
	   curl -X POST http://localhost:8080/api/rooms/create  -H "Content-Type: application/json" -d '{"name": "test-room", "created_by":"HnJO@geSAD"}'
	*/

	var roomData RoomDataRoute
	if err := parseRequestBody(r, &roomData); err != nil {
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

	userSQL := "SELECT id, name FROM users WHERE id = ?"
	userRows, err := db.Query(r.Context(), userSQL, roomData.CreatedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error querying user: %v", err)
		return
	}

	var userName string
	if userRows.Next() {
		var userID string
		if err := userRows.Scan(&userID, &userName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error scanning user: %v", err)
			return
		}
	} else {
		http.Error(w, "User not found", http.StatusBadRequest)
		log.Printf("User not found")
		return
	}

	roomID, err := GenerateID(r.Context(), db, "rooms")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error generating ID: %v", err)
		return
	}

	room := RoomData{
		ID:        roomID,
		Name:      roomData.Name,
		CreatedBy: roomData.CreatedBy,
	}

	_, err = db.Exec(r.Context(), "INSERT INTO rooms (id, name, created_by) VALUES (?, ?, ?)", room.ID, room.Name, room.CreatedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting room: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(room.ID))
	log.Printf("Room created with ID: %s For User: %s", room.ID, userName)
}

func handleRoomJoin(w http.ResponseWriter, r *http.Request) {
	/*
	   curl -X POST http://localhost:8080/api/rooms/join  -H "Content-Type: application/json" -d '{"room_id": "8uOQlvwCUD", "player_id":"HnJO@"}'
	*/

	var playerData PlayersDataRoute
	if err := parseRequestBody(r, &playerData); err != nil {
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

	roomSQL := "SELECT id, name FROM rooms WHERE id = ?"
	userRows, err := db.Query(r.Context(), roomSQL, playerData.RoomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error querying room: %v", err)
		return
	}

	var roomName string
	if userRows.Next() {
		var userID string
		if err := userRows.Scan(&userID, &roomName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error scanning room: %v", err)
			return
		}
	} else {
		http.Error(w, "Room not found", http.StatusBadRequest)
		log.Printf("Room not found")
		return
	}

	playerID, err := GenerateID(r.Context(), db, "players")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error generating ID: %v", err)
		return
	}

	room := PlayersData{
		ID:       playerID,
		RoomID:   playerData.RoomID,
		PlayerID: playerData.PlayerID,
	}

	_, err = db.Exec(r.Context(), "INSERT INTO players (id, room_id, player_id) VALUES (?, ?, ?)", room.ID, room.RoomID, room.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting player: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(room.ID))
	log.Printf("Player joined with ID: %s For Room: %s", room.PlayerID, roomName)
}
