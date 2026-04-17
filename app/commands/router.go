package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CommandHandler struct {
	store     *store
	listStore *ListStore
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		store:     NewStore(),
		listStore: NewListStore(),
	}
}

func (ch *CommandHandler) HandleCommand(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) > 0 {

		command := strings.ToUpper(value.Array[0].Bulk)

		switch command {

		case "PING":
			return ping()

		case "ECHO":
			return echo(value)

		case "SET":
			return ch.store.set(value)

		case "GET":
			return ch.store.get(value)

		case "DELETE":
			ch.store.delete(*value)
			return resp.EncodeString("OK")

		case "EXISTS":
			exists := ch.store.exists(*value)
			if exists {
				return resp.EncodeInteger(1)
			} else {
				return resp.EncodeInteger(0)
			}

		case "RPUSH":
			return ch.listStore.rpush(value)

		case "LRANGE":
			return ch.listStore.lrange(value)

		default:
			return resp.EncodeError("ERR unknown command")
		}
	}

	return resp.EncodeError("ERR unknown command")
}
