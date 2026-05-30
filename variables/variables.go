package variables

import (
	"sync"
	"time"
)

var Storage = make(map[string]string)
var Expirations = make(map[string]time.Time)
var Mu sync.Mutex
