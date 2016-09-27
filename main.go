package main

import (
	"github.com/sayevsky/godis/server"
)

func main() {
	server.NewServer().Start(false)
}
