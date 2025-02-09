package types

import "sync"

var (
	Suspicion = false
	Su        sync.RWMutex
)
