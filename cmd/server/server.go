/*
This package will be used as the server. It will listen for incoming connections from the client. It will receive logs and save them. We can send client coming from the controlcenter to the client
*/
package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"

	"cs50-romain/pcdeploy/cmd/client"
)

type Server struct {
	addr string
	port int
	ClientComputers client.ClientComputers
	Quit chan interface{} 
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr: addr,
		port: port,
		Quit: make(chan interface{}),
	}
}

func Broadcaster() {
	
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr + ":" + strconv.Itoa(s.port))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("Starting server")
	for {

		// Clean up
		//return s.stop()
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.Quit:
				fmt.Println("Quitting server")
				ln.Close()
				return nil
			default:
				fmt.Println(err)
			}
		} else {
			go s.handleConn(conn)
		}
	}
}

func (s *Server) Stop() error {
	for _, client := range s.ClientComputers.Conns {
		s.ClientComputers.Remove(*client)
		close(client.Quit)
		client.Conn.Close()
	}
	close(s.Quit)
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	// Listen for input from the client
	// Receive updates from client
	// if the client is in the foreground position, output to shell and log
	// if client is in background position, log/save.
	client := client.NewClientComputer(conn, conn.RemoteAddr().String())
	s.ClientComputers.Add(client)
	fmt.Println("A client has successfully connected: ", client.Ip)
	input := bufio.NewScanner(client.Conn)
	for input.Scan() {
		select {
		case <-client.Quit:
			fmt.Println("Closing connection")
			return
		default:
			conn.Write([]byte{'d'})
			message := input.Text()
			if client.IsForeground() {
				fmt.Println(client.Ip + ": " + message)
			}
			client.SaveLogs(message)
		}
	}
}

func (s *Server) cancelled() bool {
	select {
	case <-s.Quit:
		return true
	default:
		return false
	}
}
