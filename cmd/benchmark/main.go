package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grunclepug/gopher-cache/pkg/storage"
)

func runBenchmark(name string, store storage.Store, iterations int) {
	data := []byte("payload")

	// Write
	start := time.Now()
	for i := range iterations {
		_ = store.Put(fmt.Sprintf("key-%d", i), data)
	}
	writeDuration := time.Since(start)

	// Read
	start = time.Now()
	for i := range iterations {
		_, _ = store.Get(fmt.Sprintf("key-%d", i))
	}
	readDuration := time.Since(start)

	// Results
	fmt.Printf("--- %s ---\n", name)
	fmt.Printf("Writes (%d): %v\n", iterations, writeDuration)
	fmt.Printf("Reads  (%d): %v\n", iterations, readDuration)
	fmt.Println()
}

func main() {
	const iterations = 1000

	// 1. Benchmark Memory
	memStore := storage.NewMemoryStore()
	runBenchmark("Memory Store", memStore, iterations)

	// 2. Benchmark Disk
	diskStore, err := storage.NewDiskStore("./bench_db")
	if err != nil {
		log.Fatal(err)
	}
	runBenchmark("Disk Store", diskStore, iterations)

	// 3. Teardown
	if err := os.RemoveAll("./bench_db"); err != nil {
		log.Printf("failed to clean up: %v", err)
	}
}
