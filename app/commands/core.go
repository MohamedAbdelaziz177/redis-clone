package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func ping() []byte {
	return []byte("+PONG\r\n")

}

func echo(value *resp.Value) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value.Array[1].Bulk), value.Array[1].Bulk))
}
