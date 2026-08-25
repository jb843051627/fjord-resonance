package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSON(request *http.Request, target any) error {
	if request.Body == nil {
		return fmt.Errorf("request body is empty")
	}
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	return nil
}

func WriteJSON(writer http.ResponseWriter, status int, value any) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(value)
}

func WriteError(writer http.ResponseWriter, status int, err error) {
	_ = WriteJSON(writer, status, ErrorResponse{Error: err.Error()})
}
