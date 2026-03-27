package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/grunclepug/gopher-cache/pkg/storage"
)

type Handler struct {
	DB      storage.Store
	Verbose bool
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if h.DB == nil {
		http.Error(w, "store not initialized", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	wantsJSON := r.Header.Get("Accept") == "application/json"

	// Bucket Logic
	if strings.HasSuffix(key, "/") {
		dataMap, err := h.DB.GetBucket(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if h.Verbose {
			log.Printf("GET BUCKET %s - Found %d items\n", key, len(dataMap))
		}

		if wantsJSON {
			h.writeJSONBucket(w, dataMap)
			return
		}

		for k, v := range dataMap {
			fmt.Fprintf(w, "%s: %s\n", k, string(v))
		}
		return
	}

	// Single Key Logic
	val, err := h.DB.Get(key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.Verbose {
		log.Printf("GET %s - Success (%d bytes)\n", key, len(val))
	}

	if wantsJSON {
		h.writeJSONValue(w, val)
		return
	}

	w.Write(val)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Ingress Validation
	if r.Header.Get("Content-Type") == "application/json" {
		if !json.Valid(body) {
			http.Error(w, "invalid json payload", http.StatusBadRequest)
			return
		}
	}

	if err := h.DB.Put(key, body); err != nil {
		if errors.Is(err, storage.ErrInvalidKey) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.Verbose {
		log.Printf("POST %s - Saved (%d bytes)\n", key, len(body))
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Ingress Validation
	if r.Header.Get("Content-Type") == "application/json" {
		if !json.Valid(body) {
			http.Error(w, "invalid json payload", http.StatusBadRequest)
			return
		}
	}

	if err := h.DB.Update(key, body); err != nil {
		if errors.Is(err, storage.ErrInvalidKey) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.Verbose {
		log.Printf("PUT %s - Saved (%d bytes)\n", key, len(body))
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if err := h.DB.Delete(key); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.Verbose {
		log.Printf("DELETE %s - Success\n", key)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Helpers for "Smart" JSON promotion
func (h *Handler) writeJSONValue(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	var jsonVal any
	if err := json.Unmarshal(data, &jsonVal); err == nil {
		json.NewEncoder(w).Encode(jsonVal)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"value": string(data)})
}

func (h *Handler) writeJSONBucket(w http.ResponseWriter, data map[string][]byte) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]any)
	for k, v := range data {
		var jsonVal any
		if err := json.Unmarshal(v, &jsonVal); err == nil {
			response[k] = jsonVal
		} else {
			response[k] = string(v)
		}
	}
	json.NewEncoder(w).Encode(response)
}
