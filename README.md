# gopher-cache

A thread-safe KV store for Go. Supports swappable RAM and Disk backends through a single interface. Use it for fast in-memory caching or simple filesystem persistence.

## Overview

`gopher-cache` is built for cases where you need a pluggable storage layer. It uses standard Go primitives (sync.RWMutex) to ensure data integrity during concurrent operations.

## Features

* **Pluggable Architecture:** Common `Store` interface for all backends.
* **MemoryStore:** High-performance `map` implementation with `RWMutex` for concurrent read/write safety.
* **DiskStore:** Filesystem-backed persistence (one file per key) with automated directory management.
* **Concurrency-Safe:** Designed to handle multiple simultaneous goroutines.
* **Namespaced Buckets:** Use trailing slashes (`/foo/`) to perform recursive GET or DELETE operations.
* **Content Negotiation:** Request `application/json` via headers to get structured data trees instead of raw bytes.
* **Ingress Validation:** Automatically rejects malformed JSON if the `Content-Type` header is set.
* **Minimal:** No external dependencies. Uses Go 1.22+ standard library routing.

## Usage

### Library Example

```go
import (
    "fmt"
    "log"
    "github.com/grunclepug/gopher-cache/pkg/storage"
)

func main() {
    // Initialize MemoryStore or DiskStore
    //db := storage.NewMemoryStore()
    db, err := storage.NewDiskStore("./data")
    if err != nil {
        log.Fatalf("failed to initialize: %v", err)
    }

    // Put: Create or Overwrite a key
    key, val := "foo", []byte("payload") // Bucket Value Example: "foo/status"
    if err := db.Put(key, val); err != nil {
        log.Fatal(err)
    }

    // Update: Modify existing key (returns storage.ErrNotFound if missing)
    newVal := []byte("updated_payload")
    if err := db.Update(key, newVal); err != nil {
        log.Fatal(err)
    }

    // Get: Retrieve data
    data, err := db.Get(key)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Retrieved: %s\n", string(data))

    // Delete: Remove a key
    if err := db.Delete(key); err != nil {
        log.Fatal(err)
    }

    // GetBucket: Retrieve all keys under a prefix
    data, err := db.GetBucket("foo/")
        if err == nil {
            for k, v := range data {
            fmt.Printf("%s: %s\n", k, string(v)) // Or unmarshal if JSON
        }
    }
}
```

### Daemon (gopher-cached)

The project includes a standalone HTTP daemon that exposes the KV store over a REST API.

#### Endpoints

* **GET /health:** Heartbeat check.
* **GET /{key}:** Get raw data.
* **GET /{key} (Accept: application/json):** Returns a JSON object. If the value is a string, it's wrapped; if it's valid JSON, it's promoted to a nested object.
* **GET /{bucket}/:** Returns a list of key: value pairs.
* **GET /{bucket}/ (Accept: application/json):** Returns a single JSON tree of the entire bucket.
* **POST/PUT /{key}:** Store data. Validates JSON if Content-Type: application/json is used.
* **DELETE /{key}:** Delete a single key.
* **DELETE /{bucket}/:** Recursively delete everything in that namespace.

#### Examples

Store a value:

`curl -X POST -d "running" http://localhost:8080/nodes/srv1`


Retrieve as JSON:

`curl -H "Accept: application/json" http://localhost:8080/nodes/srv1`

Output: `{"value": "running"}`


Get a whole bucket as a JSON tree:

`curl -H "Accept: application/json" http://localhost:8080/nodes/`

Output: `{"nodes/srv1": "running", "nodes/srv2": {"load": 0.5}}`

#### Running the Daemon
```bash
# Default: DiskStore at ./data on port 8080
./bin/gopher-cached

# MemoryStore on custom port with verbose logging
./bin/gopher-cached -mem -port 9000 -v

# Help command to see all available flags
./bin/gopher-cached --help
```

#### Systemd Service
```toml
# /etc/systemd/system/gopher-cached.service
[Unit]
Description=Gopher Cache Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/gopher-cached -port 8080 -dir /var/lib/gopher-cache -v
Restart=on-failure
User=nobody
Group=nobody
# Ensure the directory exists and is writable by the user above
# sudo mkdir -p /var/lib/gopher-cache && sudo chown nobody:nobody /var/lib/gopher-cache

[Install]
WantedBy=multi-user.target
```

### Installation

#### Requirements

* Go 1.22+ (Required for range over int support)
* Make

#### Building

Clone the repository:
`git clone https://github.com/GrunclePug/gopher-cache.git`

Compile the main application:
`make`

(Optional) Compile the daemon:
`make daemon`

(Optional) Compile the performance benchmark tool:
`make benchmark`

#### Benchmarking

To compare the latency between the Memory and Disk engines on your hardware:

```bash
make benchmark
./bin/benchmark
```

#### Maintenance

Clean build artifacts and test databases:
`make clean`

Install to system (daemon):
`sudo make install` (Defaults to /usr/local/bin)

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please fork the repository and submit a pull request.

## Author

Chad Humphries |
[Website](https://grunclepug.com/) |
[GitHub Profile](https://github.com/GrunclePug)

## Other Projects

Check out some of my other projects on GitHub: [Here](https://github.com/GrunclePug?tab=repositories)
