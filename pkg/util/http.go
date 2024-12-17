package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"asynclab.club/AsyncFunction/pkg/program"
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

func HttpRequest(req *http.Request) ([]byte, error) {
	resp, err := program.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
