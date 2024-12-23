package websocketTasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Arinji2/websockets/sqlite"
)

func CreateRoomTask(ctx context.Context, taskID string, jsonData json.RawMessage) (TaskResponse, error) {
	var data RoomData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	db, err := sqlite.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	userSQL := "SELECT id, name FROM users WHERE id = ?"
	userRows, err := db.Query(ctx, userSQL, data.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	defer userRows.Close()

	var userName string
	if userRows.Next() {
		var userID string
		if err := userRows.Scan(&userID, &userName); err != nil {
			return nil, fmt.Errorf("failed to scan user data: %w", err)
		}
	} else {
		return nil, fmt.Errorf("user not found: %s", data.CreatedBy)
	}

	roomID, err := sqlite.GenerateID(ctx, db, "rooms")
	if err != nil {
		return nil, fmt.Errorf("failed to generate room ID: %w", err)
	}

	room := RoomData{
		ID:        roomID,
		Name:      data.Name,
		CreatedBy: data.CreatedBy,
	}

	_, err = db.Exec(ctx, "INSERT INTO rooms (id, name, created_by) VALUES (?, ?, ?)",
		room.ID, room.Name, room.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to insert room: %w", err)
	}

	// Keep one log entry for operational visibility
	log.Printf("Room created successfully - ID: %s, User: %s", room.ID, userName)

	return WebsocketTask[RoomData]{
		TaskID:   taskID,
		TaskData: data,
	}, nil
}

func JoinRoomTask(ctx context.Context, taskID string, jsonData json.RawMessage) (TaskResponse, error) {
	var data PlayersData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}
	db, err := sqlite.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	roomSQL := "SELECT id, name FROM rooms WHERE id = ?"
	roomRows, err := db.Query(ctx, roomSQL, data.RoomID)
	if err != nil {
		return nil, fmt.Errorf("failed to query room: %w", err)
	}
	defer roomRows.Close()

	var roomName string
	if roomRows.Next() {
		var roomID string
		if err := roomRows.Scan(&roomID, &roomName); err != nil {
			return nil, fmt.Errorf("failed to scan room data: %w", err)
		}
	} else {
		return nil, fmt.Errorf("room not found")
	}

	playerID, err := sqlite.GenerateID(ctx, db, "players")
	if err != nil {
		return nil, fmt.Errorf("failed to generate player ID: %w", err)
	}

	player := PlayersData{
		ID:       playerID,
		RoomID:   data.RoomID,
		PlayerID: data.PlayerID,
	}

	_, err = db.Exec(ctx, "INSERT INTO players (id, room_id, player_id) VALUES (?, ?, ?)",
		player.ID, player.RoomID, player.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert player: %w", err)
	}

	log.Printf("Player joined successfully - ID: %s, Room: %s", player.PlayerID, roomName)

	return WebsocketTask[PlayersData]{
		TaskID:   taskID,
		TaskData: data,
	}, nil
}

func DeleteRoomTask(ctx context.Context, taskID string, jsonData json.RawMessage) (TaskResponse, error) {
	var data PlayersData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}
	db, err := sqlite.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	dataSQL := `
		SELECT r.id, u.name 
		FROM rooms r 
		JOIN players p 
			ON r.id = p.room_id AND p.player_id = ? 
		JOIN users u
			ON p.player_id = u.id
		WHERE r.id = ?
	`
	dataRows, err := db.Query(ctx, dataSQL, data.PlayerID, data.RoomID)
	if err != nil {
		return nil, fmt.Errorf("failed to query player data: %w", err)
	}
	defer dataRows.Close()

	var playerName string
	if dataRows.Next() {
		var roomID string
		if err := dataRows.Scan(&roomID, &playerName); err != nil {
			return nil, fmt.Errorf("failed to scan player data: %w", err)
		}
	} else {
		return nil, fmt.Errorf("player not found in room")
	}

	_, err = db.Exec(ctx, "DELETE FROM players WHERE player_id = ? AND room_id = ?",
		data.PlayerID, data.RoomID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete player: %w", err)
	}

	log.Printf("Player left successfully - PlayerID: %s, Room: %s, Name: %s",
		data.PlayerID, data.RoomID, playerName)

	return WebsocketTask[PlayersData]{
		TaskID:   taskID,
		TaskData: data,
	}, nil
}
