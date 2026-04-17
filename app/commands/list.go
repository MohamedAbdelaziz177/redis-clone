package commands

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ListStore struct {
	lists map[string][]string
	mu    *sync.RWMutex
}

func NewListStore() *ListStore {
	return &ListStore{
		lists: make(map[string][]string),
		mu:    &sync.RWMutex{},
	}
}

func (ls *ListStore) rpush(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 3 {

		listName := value.Array[1].Bulk

		var elements = make([]string, 0)

		for _, arg := range value.Array[2:] {
			elements = append(elements, arg.Bulk)
		}

		ls.mu.Lock()
		defer ls.mu.Unlock()

		list, ok := ls.lists[listName]
		if !ok {
			list = []string{}
		}
		list = append(list, elements...)
		ls.lists[listName] = list

		return []byte(fmt.Sprintf(":%d\r\n", len(list)))
	}

	return []byte("-ERR invalid command format\r\n")
}

func (ls *ListStore) lrange(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) == 4 {

		listName := value.Array[1].Bulk
		start, _ := strconv.Atoi(value.Array[2].Bulk)
		end, _ := strconv.Atoi(value.Array[3].Bulk)

		ls.mu.RLock()
		defer ls.mu.RUnlock()

		list, ok := ls.lists[listName]
		if !ok {
			return resp.EncodeArray([]string{})
		}

		if start < 0 {
			start = len(list) + start
		}

		if end < 0 {
			end = len(list) + end
		}

		if end >= len(list) {
			end = len(list) - 1
		}

		if start > end || start >= len(list) || end < 0 || start < 0 {
			return resp.EncodeArray([]string{})
		}

		return resp.EncodeArray(list[start : end+1])
	}
	return resp.EncodeError("Err")
}
