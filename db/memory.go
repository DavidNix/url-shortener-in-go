package db

import (
	"sync"
)

type InMemory struct {
	sync.RWMutex
	mapping map[string]string
	tracks map[string]int64
}

func NewMemory() *InMemory {
	return &InMemory{
		mapping: make(map[string]string),
		tracks: make(map[string]int64),
	}
}

func (db *InMemory) Save(key, target string) {
	db.Lock()
	db.mapping[key] = target
	db.Unlock()
}

func (db *InMemory) Find(key string) (string, bool) {
	db.RLock()
	defer db.RUnlock()
	k, ok := db.mapping[key]
	return k, ok
}

func (db *InMemory) Track(key string) {
	db.Lock()
	defer db.Unlock()
	db.tracks[key]++
}

func (db *InMemory) Visits(key string) int64 {
	db.RLock()
	defer db.RUnlock()
	return db.tracks[key]
}
