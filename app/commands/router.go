package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CommandHandler struct {
	store *store
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		store: NewStore(),
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
			return ch.store.HSET(*value)

		case "GET":
			return ch.store.HGET(*value)

		case "DELETE":
			ch.store.HDEL(*value)
			return []byte("+OK\r\n")

		case "EXISTS":
			exists := ch.store.HEXISTS(*value)
			if exists {
				return []byte(":1\r\n")
			} else {
				return []byte(":0\r\n")
			}

		default:
			return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command))
		}
	}

	return []byte("-ERR invalid command format\r\n")
}
