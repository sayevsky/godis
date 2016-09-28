package server

import "net"
import "log"
import (
	"bufio"
	"github.com/sayevsky/godis/internal"
)

type Server struct {
	Listener     net.Listener
	dbChannel    chan interface{}
	WithEviction bool
	poisonPill   chan bool
}

func NewServer() Server {
	port := "6380"
	log.Println("Launching godis on port " + port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error starting server on "+port, err.Error())
	}

	// channel to communicate with kv-storage
	dbChannel := make(chan interface{})

	return Server{listener, dbChannel, true, make(chan bool)}

}

// connections
func (s Server) serveConnections() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Println(err.Error())
			break
		} else {
			go handle(conn, s.dbChannel)
		}
	}
}

func (s Server) Start(background bool) {
	// works with storage directly
	go ProcessCommands(s.dbChannel)

	if s.WithEviction {
		// periodically send evict message to processCommands
		go sendEvictMessages(s.dbChannel, s.poisonPill)
	}
	if background {
		go s.serveConnections()
	} else {
		s.serveConnections()
	}
}

func (s Server) Stop() {
	// will stop connections loop
	s.Listener.Close()
	// exit eviction routine
	s.poisonPill <- true
	// exit ProcessCommands
	s.dbChannel <- &internal.Quit{}
}

func handle(conn net.Conn, dbChannel chan interface{}) {
	reader := bufio.NewReader(conn)
	for {
		// handle SIGINT
		signal, err := reader.Peek(1)
		if err != nil {
			log.Println("can't peek a byte", err)
			break
		}
		if signal[0] == byte(255) {
			log.Println("Exit signal, close connection.")
			conn.Close()
			return
		}

		command, err := internal.ParseCommand(reader)

		var response internal.Response

		if err != nil {
			// if fail to parse command, send it as result
			response = internal.Response{nil, err}
			conn.Write(response.Serialize())
		} else {

			dbChannel <- command

			// send reply only for sync commands
			if ! command.GetBaseCommand().IsAsync {
				response = <-command.GetBaseCommand().ChannelWithResult
				conn.Write(response.Serialize())
			}
		}

	}
}
