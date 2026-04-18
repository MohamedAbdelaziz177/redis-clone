package commands

import (
	"strconv"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type zsetStore struct {
	zsets map[string]zsetItem
	mu    *sync.RWMutex
}

type zsetItem map[string]float64

func NewZsetStore() *zsetStore {
	return &zsetStore{
		zsets: make(map[string]zsetItem),
		mu:    &sync.RWMutex{},
	}
}

func (store *zsetStore) zadd(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) >= 4 {
		setName := value.Array[1].Bulk

		store.mu.Lock()
		defer store.mu.Unlock()

		zset, ok := store.zsets[setName]
		if !ok {
			zset = make(zsetItem)
		}

		var count int = 0
		for i := 2; i < len(value.Array); i += 2 {

			score, err := strconv.ParseFloat(value.Array[i].Bulk, 64)
			if err != nil {
				return resp.EncodeError("ERR value is not a valid float")
			}

			member := value.Array[i+1].Bulk

			score, exists := zset[member]
			if !exists {
				count++
			}
			zset[member] = score
		}

		store.zsets[setName] = zset

		return resp.EncodeInteger(count)
	}
	return resp.EncodeError("ERR Invalid resp format")
}
