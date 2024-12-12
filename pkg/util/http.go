package util

import (
	"encoding/json"
	"net/http"
)

func HttpResponse(w http.ResponseWriter, status int, message any) {
	h := w.Header()

	h.Del("Content-Length")
	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(status)
	jsonResponse, _ := json.Marshal(map[string]any{"message": message})
	w.Write([]byte(jsonResponse))
}
