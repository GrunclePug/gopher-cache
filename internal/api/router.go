package api

import (
	"net/http"
)

// NewRouter wires up the routes to the provided Handler.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// Infrastructure
	mux.HandleFunc("GET /health", h.HandleHealth)

	// Key/Bucket Operations
	mux.HandleFunc("GET /{key...}", h.HandleGet)
	mux.HandleFunc("POST /{key...}", h.HandlePost)
	mux.HandleFunc("PUT /{key...}", h.HandlePut)
	mux.HandleFunc("DELETE /{key...}", h.HandleDelete)

	return mux
}
