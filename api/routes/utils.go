package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/Arinji2/websockets/sqlite"
)

// ParseRequestBody parses the JSON body of a request into the provided struct.
func parseRequestBody(r *http.Request, data interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content type is not application/json")
	}
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is empty")
		}
		return errors.New("invalid JSON format")
	}

	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789@#$!"

// generates a 10 digit ID
func GenerateID(ctx context.Context, con *sqlite.Connection, table string) (string, error) {
	for {
		id := generateRandomString(10)
		exists, err := checkIDExists(ctx, con, id, table)
		if err != nil {
			return "", err
		}

		if !exists {
			return id, nil
		}
	}
}

func checkIDExists(ctx context.Context, con *sqlite.Connection, id string, table string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE id = ?", table)
	rows, err := con.Query(ctx, query, id)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		var count int
		if err := rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, nil
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
