package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func HandleCommand(value *resp.Value) []byte {
	if value.Type == resp.ARRAY && len(value.Array) > 0 {
		command := strings.ToUpper(value.Array[0].Bulk)
		switch command {
		case "PING":
			return ping()
		case "ECHO":
			return echo(value)
		default:
			return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command))
		}
	}
	return []byte("-ERR invalid command format\r\n")
}

func ping() []byte {
	return []byte("+PONG\r\n")

}

func echo(value *resp.Value) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value.Array[1].Str), value.Array[1].Str))
}
