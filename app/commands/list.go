package commands

import (
	"fmt"
	"strconv"
	"sync"
	"time"

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

func (ls *ListStore) lpush(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) >= 3 {

		var elements = make([]string, 0)

		for _, arg := range value.Array[2:] {
			elements = append([]string{arg.Bulk}, elements...)
		}

		listName := value.Array[1].Bulk

		ls.mu.Lock()
		defer ls.mu.Unlock()

		list, ok := ls.lists[listName]
		if !ok {
			list = []string{}
		}

		list = append(elements, list...)
		ls.lists[listName] = list

		return resp.EncodeInteger(len(list))
	}

	return resp.EncodeError("Err Handling LPUSH! ")
}

func (ls *ListStore) lrange(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) == 4 {

		ls.mu.RLock()
		defer ls.mu.RUnlock()

		listName := value.Array[1].Bulk
		start, _ := strconv.Atoi(value.Array[2].Bulk)
		end, _ := strconv.Atoi(value.Array[3].Bulk)

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

		if start < 0 {
			start = 0
		}

		fmt.Printf("LRANGE %s %d %d\n", listName, start, end)

		if start > end || start >= len(list) {
			return resp.EncodeArray([]string{})
		}

		return resp.EncodeArray(list[start : end+1])
	}
	return resp.EncodeError("Err")
}

func (ls *ListStore) llen(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) == 2 {
		listName := value.Array[1].Bulk
		ls.mu.RLock()
		defer ls.mu.RUnlock()
		list, ok := ls.lists[listName]
		if !ok {
			return resp.EncodeInteger(0)
		}
		return resp.EncodeInteger(len(list))
	}
	return resp.EncodeError("Err")

}

func (ls *ListStore) lpop(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) >= 2 {

		listName := value.Array[1].Bulk
		ls.mu.Lock()
		defer ls.mu.Unlock()

		list, ok := ls.lists[listName]

		if !ok || len(list) == 0 {
			return resp.EncodeBulk("")
		}

		if len(value.Array) == 3 {
			count, _ := strconv.Atoi(value.Array[2].Bulk)
			if count > 0 && count < len(list) {

				elements := list[:count]
				ls.lists[listName] = list[count:]
				return resp.EncodeArray(elements)

			}
		} else {

			element := list[0]
			ls.lists[listName] = list[1:]
			return resp.EncodeBulk(element)
		}

	}
	return resp.EncodeError("Err")
}

func (ls *ListStore) blpop(value *resp.Value) []byte {

	if value.Type != resp.ARRAY || len(value.Array) < 2 {
		return resp.EncodeError("-ERR wrong number of arguments for 'blpop' command")
	}

	listName := value.Array[1].Bulk
	timeoutSec := 0.0

	if len(value.Array) >= 3 {
		t, err := strconv.ParseFloat(value.Array[2].Bulk, 64)
		if err == nil {
			timeoutSec = t
		}
	}

	var deadline time.Time
	if timeoutSec > 0 {
		deadline = time.Now().Add(time.Duration(timeoutSec * float64(time.Second)))
	}

	for {

		ls.mu.Lock()
		list, exists := ls.lists[listName]

		if exists && len(list) > 0 {
			ele := list[0]
			ls.lists[listName] = list[1:]
			ls.mu.Unlock()

			return resp.EncodeArray([]string{listName, ele})
		}

		if timeoutSec > 0 && time.Now().After(deadline) {
			ls.mu.Unlock()
			return resp.EncodeNull()
		}

		ls.mu.Unlock()
		time.Sleep(50 * time.Millisecond)

	}

}
