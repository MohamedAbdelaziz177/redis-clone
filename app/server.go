package main

import (
	"github.com/MohamedAbdelaziz177/redis-clone/app/commands"
	"github.com/MohamedAbdelaziz177/redis-clone/app/persistance"
)

type Server struct {
	port           int
	cmdHandler    *commands.CommandHandler
	aof           *persistance.AOF
}