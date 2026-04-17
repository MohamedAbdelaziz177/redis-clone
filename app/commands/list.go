package commands

import (
	"fmt"
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

		for _, arg := range value.Array[3:] {
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
