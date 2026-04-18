package commands

import (
	"strconv"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type zsetStore struct {
	zsets map[string][]zsetItem
	mu    *sync.RWMutex
}

type zsetItem struct {
	member string
	score  float64
}

func NewZsetStore() *zsetStore {
	return &zsetStore{
		zsets: make(map[string][]zsetItem),
		mu:    &sync.RWMutex{},
	}
}

func (zs *zsetStore) zadd(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) >= 4 {
		setName := value.Array[1].Bulk

		zs.mu.Lock()
		defer zs.mu.Unlock()

		zset, ok := zs.zsets[setName]
		if !ok {
			zset = []zsetItem{}
		}

		for i := 2; i < len(value.Array); i += 2 {

			score, err := strconv.ParseFloat(value.Array[i].Bulk, 64)
			if err != nil {
				return resp.EncodeError("ERR value is not a valid float")
			}

			member := value.Array[i+1].Bulk

			item := zsetItem{
				member: member,
				score:  score,
			}

			zset = append(zset, item)
		}

		zs.zsets[setName] = zset

		return resp.EncodeInteger(len(zset))
	}
	return resp.EncodeError("ERR Invalid resp format")
}
