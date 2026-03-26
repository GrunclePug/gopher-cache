package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/grunclepug/gopher-cache/pkg/storage"
)

var Version = "dev"

func main() {
	// Flags
	port := flag.Int("port", 8080, "TCP port to listen on")
	useMemory := flag.Bool("mem", false, "Use in-memory storage instead of disk (non-persistent)")
	dataDir := flag.String("dir", "./data", "Directory for DiskStore persistence")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	version := flag.Bool("version", false, "print version")

	flag.Parse()

	if *version {
		fmt.Printf("gopher-cache %s\n", Version)
		os.Exit(0)
	}

	// Initialize Store based on flags
	var db storage.Store
	var err error

	if *useMemory {
		if *verbose {
			log.Println("Initializing MemoryStore (Volatile)")
		}
		db = storage.NewMemoryStore()
	} else {
		if *verbose {
			log.Printf("Initializing DiskStore at %s (Persistent)\n", *dataDir)
		}
		db, err = storage.NewDiskStore(*dataDir)
		if err != nil {
			log.Fatalf("failed to init disk store: %v", err)
		}
	}

	mux := http.NewServeMux()

	// GET /health -> Heartbeat check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "store not initialized", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// GET /{key} -> Retrieve data
	mux.HandleFunc("GET /{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		val, err := db.Get(key)
		if err != nil {
			if *verbose {
				log.Printf("GET %s - Not Found\n", key)
			}
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if *verbose {
			log.Printf("GET %s - Success (%d bytes)\n", key, len(val))
		}
		w.Write(val)
	})

	// POST /{key} -> Put data
	mux.HandleFunc("POST /{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if err := db.Put(key, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if *verbose {
			log.Printf("POST %s - Saved (%d bytes)\n", key, len(body))
		}
		w.WriteHeader(http.StatusCreated)
	})

	// PUT /{key} -> Update existing data
	mux.HandleFunc("PUT /{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if err := db.Update(key, body); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if *verbose {
			log.Printf("PUT %s - Updated (%d bytes)\n", key, len(body))
		}
		w.WriteHeader(http.StatusOK)
	})

	// DELETE /{key} -> Remove data
	mux.HandleFunc("DELETE /{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if err := db.Delete(key); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if *verbose {
			log.Printf("DELETE %s - Success\n", key)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("gopher-cached listening on %s (Memory Mode: %v)\n", addr, *useMemory)
	log.Fatal(http.ListenAndServe(addr, mux))
}
