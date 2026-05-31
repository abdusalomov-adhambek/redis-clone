package persistence

import (
	"encoding/json"
	"goredisclone/variables"
	"log"
	"os"
)

// save saves the current state of the storage to disk.
func Save() {
	variables.Mu.Lock()
	defer variables.Mu.Unlock()

	var db = variables.DB{
		Expirations: variables.Expirations,
		Storage:     variables.Storage,
	}

	data, err := json.Marshal(db)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("./database.json", data, 0644); err != nil {
		panic(err)
	}
}

// load loads the state of the storage from disk.
func Load() {

	var db variables.DB

	readFile, err := os.ReadFile("./database.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("database.json not found, starting with empty database")
			return
		}
		panic(err)
	}

	if err := json.Unmarshal(readFile, &db); err != nil {
		panic(err)
	}

	variables.Mu.Lock()

	variables.Storage = db.Storage
	variables.Expirations = db.Expirations

	variables.Mu.Unlock()

}
