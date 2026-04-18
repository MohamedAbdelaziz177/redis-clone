package commands

import (
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type hashStore struct {
	hashmap map[string]map[string]interface{}
	mu      *sync.RWMutex
}

func NewHashStore() *hashStore {
	return &hashStore{
		hashmap: make(map[string]map[string]interface{}),
		mu:      &sync.RWMutex{},
	}
}

func (hs *hashStore) hset(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 4 {

		mpKey := value.Array[1].Bulk

		hs.mu.Lock()
		defer hs.mu.Unlock()

		mp, ok := hs.hashmap[mpKey]

		if !ok {
			mp = make(map[string]interface{})
		}

		c := 0
		for i := 2; i < len(value.Array)-1; i += 2 {
			field := value.Array[i].Bulk
			val := value.Array[i+1].Bulk
			mp[field] = val
			c++
		}

		hs.hashmap[mpKey] = mp

		return resp.EncodeInteger(c)
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (hs *hashStore) hget(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) == 3 {

		mpKey := value.Array[1].Bulk
		field := value.Array[2].Bulk

		hs.mu.RLock()
		defer hs.mu.RUnlock()

		mp, ok := hs.hashmap[mpKey]
		if !ok {
			return resp.EncodeNull()
		}

		val, exists := mp[field]
		if !exists {
			return resp.EncodeNull()
		}

		return resp.EncodeBulk(val.(string))
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (hs *hashStore) hdel(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) == 3 {

		mpKey := value.Array[1].Bulk
		field := value.Array[2].Bulk

		hs.mu.Lock()
		defer hs.mu.Unlock()

		mp, ok := hs.hashmap[mpKey]
		if !ok {
			return resp.EncodeError("No Such Hash_Key")
		}

		_, exists := mp[field]
		if !exists {
			return resp.EncodeError("Field does not exist")
		}

		delete(mp, field)
		return resp.EncodeInteger(1)
	}
	return resp.EncodeError("Err Invalid resp format")
}

func (hs *hashStore) hexists(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) == 3 {

		mpKey := value.Array[1].Bulk
		field := value.Array[2].Bulk

		hs.mu.RLock()
		defer hs.mu.RUnlock()

		mp, ok := hs.hashmap[mpKey]
		if !ok {
			return resp.EncodeError("No Such Hash_Key")
		}

		_, exists := mp[field]
		if !exists {
			return resp.EncodeInteger(0)
		}

		return resp.EncodeInteger(1)
	}
	return resp.EncodeError("Err Invalid resp format")
}
