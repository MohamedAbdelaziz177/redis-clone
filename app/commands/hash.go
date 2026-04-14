package commands

import "sync"

type store struct {
	hashmap map[string]string
	mu      *sync.RWMutex
}

func NewStore() *store {
	return &store{
		hashmap: make(map[string]string),
		mu:      &sync.RWMutex{},
	}
}

func (s *store) hget(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.hashmap[key]
	return val, ok
}

func (s *store) hset(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hashmap[key] = value
}

func (s *store) hgetall() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.hashmap
}

func (s *store) hexists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.hashmap[key]
	return ok
}

func (s *store) hdel(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.hashmap, key)
}
