package sqlite

import (
	"context"
	"fmt"

	"github.com/Arinji2/websockets/types"
)

func AuthenticateUser(userData types.UserData) error {
	authSQL := `SELECT id FROM users where email = ? AND session_id = ?`
	db, err := NewConnection()
	if err != nil {
		return fmt.Errorf("error creating database connection: %v", err)
	}
	defer db.Close()

	var userID string
	rows, err := db.Query(context.Background(), authSQL, userData.Email, userData.SessionID)
	if err != nil {
		return fmt.Errorf("error authenticating user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&userID)
		if err != nil {
			return fmt.Errorf("error scanning user ID: %v", err)
		}
	}
	if userID == "" {
		return fmt.Errorf("user not found")
	}
	return nil
}
