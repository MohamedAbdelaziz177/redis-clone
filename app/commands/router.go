package commands

import (
	"fmt"
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
			return []byte("+OK\r\n")

		case "EXISTS":
			exists := ch.store.exists(*value)
			if exists {
				return []byte(":1\r\n")
			} else {
				return []byte(":0\r\n")
			}

		case "RPUSH":
			return ch.listStore.rpush(value)

		default:
			return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command))
		}
	}

	return []byte("-ERR invalid command format\r\n")
}
