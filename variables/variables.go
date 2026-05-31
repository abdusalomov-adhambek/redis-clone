package variables

import (
	"sync"
	"time"
)

var Storage = make(map[string]string)
var Expirations = make(map[string]time.Time)
var Mu sync.Mutex

type DB struct {
	Storage     map[string]string    `json:"storage"`
	Expirations map[string]time.Time `json:"expirations"`
}
