package commands

import (
	"sync"

	"github.com/MohamedAbdelaziz177/redis-clone/app/resp"
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

func (hs *hashStore) hset(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 4 {

		if (len(value.Array)-2)%2 != 0 {
			return resp.EncodeError("ERR syntax error")
		}

		mpKey := value.Array[1].Bulk

		hs.mu.Lock()
		defer hs.mu.Unlock()

		mp, ok := hs.hashmap[mpKey]

		if !ok {
			mp = make(map[string]string)
		}

		c := 0
		for i := 2; i < len(value.Array); i += 2 {
			field := value.Array[i].Bulk
			val := value.Array[i+1].Bulk

			_, exists := mp[field]
			if !exists {
				c++
			}

			mp[field] = val
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

		return resp.EncodeBulk(val)
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (hs *hashStore) hdel(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 3 {

		mpKey := value.Array[1].Bulk

		hs.mu.Lock()
		defer hs.mu.Unlock()

		mp, ok := hs.hashmap[mpKey]
		if !ok {
			return resp.EncodeInteger(0)
		}

		c := 0
		for i := 2; i < len(value.Array); i++ {

			_, exists := mp[value.Array[i].Bulk]
			if exists {
				c++
				delete(mp, value.Array[i].Bulk)
			}
		}

		return resp.EncodeInteger(c)
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
			return resp.EncodeInteger(0)
		}

		_, exists := mp[field]
		if !exists {
			return resp.EncodeInteger(0)
		}

		return resp.EncodeInteger(1)
	}
	return resp.EncodeError("Err Invalid resp format")
}
