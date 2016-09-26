package server

import "net"
import "log"
import (
	"bufio"
	"github.com/sayevsky/godis/internal"
)


func Start() {

	port := "6380"
	log.Println("Launching godis on port " + port)

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatal("Error starting server on " + port, err.Error())
	}

	defer listener.Close()

	// channel to communicate with kv-storage
	dbChannel := make(chan interface{}, 10)

	go ProcessCommands(dbChannel, true)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting ", err.Error())
		} else {
			go handle(conn, dbChannel)
		}


	}

}

func handle(conn net.Conn, dbChannel chan interface{}){
	reader := bufio.NewReader(conn)
	for {
		// handle SIGINT
		signal, _ := reader.Peek(1)
		if signal[0] == byte(255){
			log.Println("Exit signal, close connection.")
			conn.Close()
			return
		}

		command, _ := internal.ParseCommand(reader)

		dbChannel <- command

		conn.Write((<- command.GetBaseCommand().ChannelWithResult).Serialize())


	}
}
