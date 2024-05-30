/*
This package will be used as the server. It will listen for incoming connections from the client. It will receive logs and save them. We can send client coming from the controlcenter to the client
*/
package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	addr string
	port int
}

type ClientComputer struct {
	connection net.Conn
	addr string
	logs []string
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr: addr,
		port: port,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.addr + ":" + strconv.Itoa(s.port))
	if err != nil {
		return err
	}

	fmt.Println("Starting server")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// Listen for input from the client
	// Receive updates from client
	// if the client is in the foreground position, output to shell and log
	// if client is in background position, log/save.
	client := &ClientComputer{
		connection: conn,
		addr: conn.RemoteAddr().String(),
		logs: []string{},
	}
	input := bufio.NewScanner(client.connection)
	for input.Scan() {
		fmt.Println(client.addr + ": " + input.Text())
	}
}
