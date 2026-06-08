// Package persistence handles reading and writing the server state to disk so
// that data survives process restarts.
package persistence

import (
	"encoding/json"
	"goredisclone/variables"
	"log"
	"os"
)

// Save serialises the current in-memory state (Storage and Expirations) into
// JSON and writes it atomically to database.json. It holds the mutex for the
// duration of the marshal so the snapshot is consistent. Called after every
// successful SET so no writes are lost on restart.
func Save() {
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	// Snapshot both maps into the serialisable DB struct.
	var db = variables.DB{
		Expirations: variables.Expirations,
		Storage:     variables.Storage,
	}

	// Marshal the snapshot to JSON bytes.
	data, err := json.Marshal(db)
	if err != nil {
		log.Println("database.json marshal error: ", err)
		return
	}

	// Write the JSON to disk, creating or truncating the file as needed.
	if err := os.WriteFile("./database.json", data, 0644); err != nil {
		log.Println("database.json write file error: ", err)
	}
}

// Load reads database.json from disk and restores Storage and Expirations
// into memory. It is called once at startup before the server begins
// accepting connections. If the file does not exist the server starts with
// an empty dataset, which is the normal first-run behaviour.
func Load() {

	var db variables.DB

	// Read the raw JSON bytes from disk.
	readFile, err := os.ReadFile("./database.json")
	if err != nil {
		if os.IsNotExist(err) {
			// First run: no persistence file yet; start fresh.
			log.Println("database.json not found, starting with empty database")
			return
		}
		log.Println("database.json read file error: ", err)
	}

	// Deserialise the JSON into the DB snapshot struct.
	if err := json.Unmarshal(readFile, &db); err != nil {
		log.Println("database.json unmarshal error: ", err)
		return
	}

	// Replace the in-memory maps with the restored data under the mutex.
	variables.Mu.Lock()

	variables.Storage = db.Storage
	variables.Expirations = db.Expirations

	variables.Mu.Unlock()
	log.Println("data loaded successfully")
}
