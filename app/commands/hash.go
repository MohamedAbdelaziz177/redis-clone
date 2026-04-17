package commands

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

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

func (s *store) HGET(value *resp.Value) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value.Type == resp.ARRAY && len(value.Array) == 2 {

		key := value.Array[1].Bulk
		val, _ := s.hashmap[key]

		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
	}

	return []byte(fmt.Sprintf("$%d\r\n", -1))
}

func (s *store) HSET(value *resp.Value) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value.Type == resp.ARRAY && len(value.Array) == 3 {
		key := value.Array[1].Bulk
		val := value.Array[2].Bulk

		s.hashmap[key] = val
		return []byte("+OK\r\n")
	}
	return []byte("-ERR invalid command format\r\n")
}

func (s *store) HGETALL() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.hashmap
}

func (s *store) HEXISTS(value resp.Value) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value.Type == resp.ARRAY && len(value.Array) == 2 {

		key := value.Array[1].Bulk
		_, ok := s.hashmap[key]

		return ok
	}
	return false
}

func (s *store) HDEL(value resp.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value.Type == resp.ARRAY && len(value.Array) == 2 {

		key := value.Array[1].Bulk

		_, ok := s.hashmap[key]
		if ok {
			delete(s.hashmap, key)
		}
	}
}
