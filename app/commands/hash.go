package commands

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ValueItem struct {
	value      string
	expiration int64
}

type store struct {
	hashmap map[string]ValueItem
	mu      *sync.RWMutex
}

func NewStore() *store {
	return &store{
		hashmap: make(map[string]ValueItem),
		mu:      &sync.RWMutex{},
	}
}

func (s *store) HGET(value *resp.Value) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value.Type == resp.ARRAY && len(value.Array) == 2 {

		key := value.Array[1].Bulk
		val, _ := s.hashmap[key]

		if val.expiration != 0 && time.Now().UnixMilli() > val.expiration {
			delete(s.hashmap, key)
			return []byte(fmt.Sprintf("$%d\r\n", -1))
		}

		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val.value), val.value))
	}

	return []byte(fmt.Sprintf("$%d\r\n", -1))
}

func (s *store) HSET(value *resp.Value) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value.Type == resp.ARRAY && len(value.Array) >= 3 {
		key := value.Array[1].Bulk
		val := value.Array[2].Bulk

		item := ValueItem{
			value:      val,
			expiration: int64(0),
		}

		if len(value.Array) == 5 && strings.ToUpper(value.Array[3].Bulk) == "PX" {
			if expirationMs, err := strconv.Atoi(value.Array[4].Bulk); err == nil {
				item.expiration = time.Now().UnixMilli() + int64(expirationMs)
			}
		}

		s.hashmap[key] = item

		return []byte("+OK\r\n")
	}
	return []byte("-ERR invalid command format\r\n")
}

func (s *store) HGETALL() map[string]ValueItem {
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
