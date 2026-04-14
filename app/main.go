package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {

		parser := NewParser(bufio.NewReader(conn))

		value, err := parser.ParseValue()
		if err != nil {
			fmt.Println("Error parsing value: ", err.Error())
			return
		}

		if value.Type == ARRAY && len(value.Array) > 0 && value.Array[0].Str == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if value.Type == STRING && value.Str == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if value.Type == BULK && value.Bulk == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else if value.Type == INTEGER && value.Int == 0 { // This is just an example, you can adjust the condition as needed
			conn.Write([]byte("+PONG\r\n"))
		} else if value.Type == ERROR && value.Err == "PING" {
			conn.Write([]byte("+PONG\r\n"))
		} else {
			conn.Write([]byte("-ERR unknown command\r\n"))
		}

	}

}

func parseArray(conn net.Conn) []string {

	rd := bufio.NewReader(conn)

	rd.ReadByte()
	str, err := rd.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading string: ", err.Error())
		return nil
	}

	str = strings.TrimRight(str, "\r\n")
	str = strings.TrimLeft(str, "*")

	count := 0
	fmt.Sscanf(str, "%d", &count)

	result := make([]string, count)

	for i := 0; i < count; i++ {
		str, _ = rd.ReadString('\n')
		str = strings.TrimRight(str, "\r\n")
		str = strings.TrimLeft(str, "$")

		var length int
		fmt.Sscanf(str, "%d", &length)

		str, _ = rd.ReadString('\n')
		str = strings.TrimRight(str, "\r\n")

		result[i] = str
	}

	return result
}

// +PING\r\n

func handleConnection(conn net.Conn) {
	defer conn.Close()

	rd := bufio.NewReader(conn)

	rd.ReadByte()
	str, err := rd.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading string: ", err.Error())
		return
	}

	str = strings.TrimRight(str, "\r\n")

	if str == "PING" {
		conn.Write([]byte("+PONG\r\n"))
	}
}
