package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grunclepug/gopher-cache/internal/api"
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

	// Initialize API Handlers and Router
	handler := &api.Handler{
		DB:      db,
		Verbose: *verbose,
	}
	router := api.NewRouter(handler)

	// Server Setup
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("gopher-cached %s listening on %s (Memory Mode: %v)\n", Version, addr, *useMemory)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
