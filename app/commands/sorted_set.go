package commands

import (
	"sort"
	"strconv"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type zsetStore struct {
	zsets map[string]zsetItem
	mu    *sync.RWMutex
}

type zsetItem map[string]float64

type ZEntry struct {
	Member string
	Score  float64
}

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

func (store *zsetStore) zrank(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) == 3 {
		setName := value.Array[1].Bulk
		member := value.Array[2].Bulk

		store.mu.RLock()
		defer store.mu.RUnlock()

		sortedList := store.getSorted(setName)

		if sortedList == nil {
			return resp.EncodeBulk("")
		}

		for i, entry := range sortedList {
			if entry.Member == member {
				return resp.EncodeInteger(i)
			}
		}

		return resp.EncodeBulk("")
	}

	return resp.EncodeError("ERR Invalid resp format")
}

func (z *zsetStore) getSorted(key string) []ZEntry {

	item, ok := z.zsets[key]
	if !ok {
		return nil
	}

	entries := make([]ZEntry, 0, len(item))
	for member, score := range item {
		entries = append(entries, ZEntry{member, score})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		return entries[i].Member > entries[j].Member
	})

	return entries
}
