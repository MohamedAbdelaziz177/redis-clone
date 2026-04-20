package commands

import (
	"sync"

	"github.com/MohamedAbdelaziz177/redis-clone/app/resp"
)

type set map[string]struct{}

type SetStore struct {
	sets map[string]set
	mu   *sync.RWMutex
}

func NewSetStore() *SetStore {
	return &SetStore{
		sets: make(map[string]set),
		mu:   &sync.RWMutex{},
	}
}

func (ss *SetStore) sadd(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 3 {
		setName := value.Array[1].Bulk

		ss.mu.Lock()
		defer ss.mu.Unlock()

		sett, ok := ss.sets[setName]
		if !ok {
			sett = make(set)
		}

		for _, arg := range value.Array[2:] {

			if _, exists := sett[arg.Bulk]; !exists {
				sett[arg.Bulk] = struct{}{}
			}
		}

		ss.sets[setName] = sett

		return resp.EncodeInteger(len(sett))
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (ss *SetStore) srem(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 3 {
		setName := value.Array[1].Bulk

		ss.mu.Lock()
		defer ss.mu.Unlock()

		sett, ok := ss.sets[setName]
		if !ok {
			return resp.EncodeInteger(0)
		}

		c := 0
		for _, arg := range value.Array[2:] {
			if _, exists := sett[arg.Bulk]; exists {
				delete(sett, arg.Bulk)
				c++
			}
		}

		ss.sets[setName] = sett

		return resp.EncodeInteger(c)
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (ss *SetStore) smembers(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) == 2 {

		setName := value.Array[1].Bulk

		ss.mu.RLock()
		defer ss.mu.RUnlock()

		sett, ok := ss.sets[setName]
		if !ok {
			return resp.EncodeArray([]string{})
		}

		var members []string

		for mem, _ := range sett {
			members = append(members, mem)
		}

		return resp.EncodeArray(members)
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (ss *SetStore) sismember(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) == 3 {
		setName := value.Array[1].Bulk
		member := value.Array[2].Bulk

		ss.mu.RLock()
		defer ss.mu.RUnlock()

		sett, ok := ss.sets[setName]
		if !ok {
			return resp.EncodeInteger(0)
		}

		if _, exists := sett[member]; exists {
			return resp.EncodeInteger(1)
		}

		return resp.EncodeInteger(0)
	}

	return resp.EncodeError("Err Invalid resp format")
}

func (ss *SetStore) scard(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) == 2 {
		setName := value.Array[1].Bulk

		ss.mu.RLock()
		defer ss.mu.RUnlock()

		sett, ok := ss.sets[setName]
		if !ok {
			return resp.EncodeInteger(0)
		}

		return resp.EncodeInteger(len(sett))
	}

	return resp.EncodeError("Err Invalid resp format")
}
