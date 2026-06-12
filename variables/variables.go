// Package variables holds shared in-memory state used across the server:
// the key-value store, per-key expiration timestamps, and a mutex that
// guards concurrent access to both maps.
package variables

import (
	"net"
	"sync"
	"time"
)

// Storage is the main in-memory key-value store.
var Storage = make(map[string]string)

// Expirations maps each key to the absolute time at which it expires.
// Keys without an entry here never expire.
var Expirations = make(map[string]time.Time)

// Password is the server password clients must supply via AUTH.
var Password = "1"

// AuthenticatedClients tracks which active connections have passed AUTH.
var AuthenticatedClients = map[net.Conn]bool{}

// Mu protects Storage and Expirations from concurrent read/write races.
var Mu sync.Mutex

// DB is the serialisable snapshot of the server state written to and read
// from the JSON persistence file.
type DB struct {
	Storage     map[string]string    `json:"storage"`
	Expirations map[string]time.Time `json:"expirations"`
}
