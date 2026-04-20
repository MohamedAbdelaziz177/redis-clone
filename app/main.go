package main

import (
	"fmt"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	server, err := initServer(":6379", "./persistance/redis.conf")
	if err != nil {
		fmt.Println("Error - Server Couldn't Start: ", err)
	}
	server.listenAndServer()
}
