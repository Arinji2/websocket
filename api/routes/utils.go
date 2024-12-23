package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
