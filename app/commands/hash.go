package commands

import (
	"sync"
)

type hashStore struct {
	hashmap map[string]map[string]string
	mu      *sync.RWMutex
}

func NewHashStore() *hashStore {
	return &hashStore{
		hashmap: make(map[string]map[string]string),
		mu:      &sync.RWMutex{},
	}
}
