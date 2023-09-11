package httpx

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes v as a JSON response with the given status.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

func Error(w http.ResponseWriter, status int, err error) {
	JSON(w, status, map[string]string{"error": err.Error()})
}
