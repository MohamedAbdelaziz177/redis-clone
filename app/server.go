package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/MohamedAbdelaziz177/redis-clone/app/commands"
	"github.com/MohamedAbdelaziz177/redis-clone/app/persistance"
	"github.com/MohamedAbdelaziz177/redis-clone/app/resp"
)

type Server struct {
	port       string
	aof        *persistance.AOF
	cmdHandler *commands.CommandHandler
}

var writeCmds = map[string]bool{
	"SET":    true,
	"DELETE": true,
	"HSET":   true,
	"HDEL":   true,
	"LPUSH":  true,
	"RPUSH":  true,
	"LPOP":   true,
	"BLPOP":  true,
	"SADD":   true,
	"SREM":   true,
}

func initServer(port string, aofConfigPath string) (*Server, error) {

	config, err := persistance.ReadAofConfig(aofConfigPath)

	if err != nil {
		return nil, fmt.Errorf("Error Reading AOF config file when Intiating The Server")
	}

	aof, err := persistance.NewAOF(config)
	if err != nil {
		return nil, fmt.Errorf("Error Iniating AOF Persistance when Intiating The Server")
	}

	ch := commands.NewCommandHandler()

	values, err := aof.ReadEntries()
	if err != nil {
		return nil, err
	}

	for _, v := range values {
		ch.HandleCommand(&v)
	}

	return &Server{
		port:       port,
		cmdHandler: ch,
		aof:        aof,
	}, nil
}

func (s *Server) listenAndServer() {

	l, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatal("Server is Unable to listen")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error Accepting Conn.. ")
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	parser := resp.NewParser(conn)

	for {

		value, err := parser.ParseValue()
		if err != nil {
			fmt.Println("Error parsing value: ", err.Error())
			return
		}

		conn.Write(s.cmdHandler.HandleCommand(&value))

		if _, ok := writeCmds[strings.ToUpper(value.Array[0].Bulk)]; ok {
			s.aof.Append(&value)

			if err := s.aof.Wr.Flush(); err != nil {
				return
			}
		}
	}
}
