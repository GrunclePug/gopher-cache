package main

import (
	"fmt"
	"log"

	"github.com/grunclepug/gopher-cache/pkg/storage"
)

var Version = "dev"

func main() {
	// Database Type
	db := storage.NewMemoryStore()
	//db, err := storage.NewDiskStore("./test_db")
	//if err != nil {
	//	log.Fatalf("failed to initialize disk store: %v", err)
	//}

	// Sample Data
	key := "config"
	value := []byte("v0.1.0")

	// Put
	if err := db.Put(key, value); err != nil {
		log.Printf("Put failed: %v", err)
	}

	// Get
	data, err := db.Get(key)
	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		fmt.Printf("Retrieved %s: %v\n", key, string(data))
	}

	// Update
	newValue := []byte("v0.1.1")
	if err := db.Update(key, newValue); err != nil {
		log.Printf("Update failed: %v", err)
	}

	// Get (confirm Update)
	data, err = db.Get(key)
	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		fmt.Printf("Retrieved %s: %v\n", key, string(data))
	}

	// Update non-existent key (should fail)
	err = db.Update("config", []byte("data"))
	if err == storage.ErrNotFound {
		fmt.Println("Correctly identified missing key during Update")
	}

	// Delete
	if err := db.Delete(key); err != nil {
		log.Printf("Delete failed: %v", err)
	}
}
