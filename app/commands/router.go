package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CommandHandler struct {
	store     *store
	listStore *ListStore
	setStore  *SetStore
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		store:     NewStore(),
		listStore: NewListStore(),
		setStore:  NewSetStore(),
	}
}

func (ch *CommandHandler) HandleCommand(value *resp.Value) []byte {

	if value.Type == resp.ARRAY && len(value.Array) > 0 {

		command := strings.ToUpper(value.Array[0].Bulk)

		switch command {

		case "PING":
			return ch.ping()

		case "ECHO":
			return ch.echo(value)

		case "TYPE":
			return ch.typ(value)

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

		case "LPUSH":
			return ch.listStore.lpush(value)

		case "LRANGE":
			return ch.listStore.lrange(value)

		case "LLEN":
			return ch.listStore.llen(value)

		case "LPOP":
			return ch.listStore.lpop(value)

		case "BLPOP":
			return ch.listStore.blpop(value)

		case "SADD":
			return ch.setStore.sadd(value)

		case "SREM":
			return ch.setStore.srem(value)

		case "SMEMBERS":
			return ch.setStore.smembers(value)

		case "SISMEMBER":
			return ch.setStore.sismember(value)

		case "SCARD":
			return ch.setStore.scard(value)

		default:
			return resp.EncodeError("ERR unknown command")
		}
	}

	return resp.EncodeError("ERR unknown command")
}
