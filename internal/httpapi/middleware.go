package httpapi

import (
	"fmt"
	"net/http"
	"time"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := request.Header.Get("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		writer.Header().Set("X-Request-ID", id)
		next.ServeHTTP(writer, request)
	})
}

func method(writer http.ResponseWriter, request *http.Request, expected string) bool {
	if request.Method == expected {
		return true
	}
	writer.Header().Set("Allow", expected)
	writer.WriteHeader(http.StatusMethodNotAllowed)
	return false
}
