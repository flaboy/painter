package httpserver

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

func NewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/healthz":
			writeOK(w, http.StatusOK, map[string]any{"ok": true})
		default:
			writeErr(w, http.StatusNotFound, "NOT_FOUND", "route not found")
		}
	})
}

func writeOK(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{Data: data, Error: nil})
}

func writeErr(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(envelope{
		Data: nil,
		Error: map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
