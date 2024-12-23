package websocketTasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Arinji2/websockets/sqlite"
)

func AuthenticateUser(ctx context.Context, taskID string, jsonData json.RawMessage) (TaskResponse, error) {
	var data UserData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}
	authSQL := `SELECT id FROM users where email = ? AND session_id = ?`
	db, err := sqlite.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("error creating database connection: %v", err)
	}
	defer db.Close()

	var userID string
	rows, err := db.Query(context.Background(), authSQL, data.Email, data.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error authenticating user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID: %v", err)
		}
	}
	if userID == "" {
		return nil, fmt.Errorf("user not found")
	}
	return WebsocketTask[UserData]{
		TaskID:   taskID,
		TaskData: data,
	}, nil
}
