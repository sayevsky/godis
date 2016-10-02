package main

import (
	"github.com/sayevsky/godis/server"
	"flag"
)

func main() {
	port := flag.String("port", "6380", "port to listen")
	flag.Parse()
	server.NewServerWithPort(*port).Start(false)
}
