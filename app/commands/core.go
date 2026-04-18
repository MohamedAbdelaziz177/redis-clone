package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (ch *CommandHandler) ping() []byte {
	return resp.EncodeString("PONG")

}

func (ch *CommandHandler) echo(value *resp.Value) []byte {
	return resp.EncodeBulk(value.Array[1].Bulk)
}

func (ch *CommandHandler) typ(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) == 2 {
		key := value.Array[1].Bulk

		ch.store.mu.RLock()
		defer ch.store.mu.RUnlock()
		if _, exists := ch.store.hashmap[key]; exists {
			return resp.EncodeString("string")
		}

		if _, exists := ch.listStore.lists[key]; exists {
			return resp.EncodeString("list")
		}

		if _, exists := ch.setStore.sets[key]; exists {
			return resp.EncodeString("set")
		}

		return resp.EncodeString("none")
	}

	return resp.EncodeString("none")
}
