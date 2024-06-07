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
	Quit chan bool
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr: addr,
		port: port,
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
		select {
		case <-s.Quit:
			fmt.Println("Quitting server")
			// Clean up
			return s.stop()
		default:
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
			}
			go s.handleConn(conn)
		}
	}
}

func (s *Server) stop() error {
	for _, client := range s.ClientComputers.Conns {
		s.ClientComputers.Remove(*client)
	}
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
			fmt.Println(client.Ip + ": " + message)
			client.SaveLogs(message)
		}
	}
}
