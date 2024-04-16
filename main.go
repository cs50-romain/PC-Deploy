package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	//"io"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn 	net.Conn
	ip 	string
	Status  bool		// foreground = true.
	Logs	[]string
	rcv	<-chan string
}

type Clients struct {
	Data []Client
}

func (c *Client) SaveLogs(log string) {
	c.Logs = append(c.Logs, log)
}

func (c *Clients) Remove(client Client) {
	// Remove client from the array
	for i, conns := range c.Data {
		if client.ip == conns.ip {
			c.Data = append(c.Data[:i], c.Data[i+1:]...)	
		}
	}
}

func (c *Clients) Add(client Client) {
	c.Data = append(c.Data, client)
}

var clients Clients

func main() {
	fmt.Println("STARTING CONTROL CENTER & SERVER!")

	var wg sync.WaitGroup
	var sendCmds = make(chan string)

	ln, err := net.Listen("tcp4", ":6969")
	if err != nil {
		log.Fatalf("Listen error: %s\n", err)
	}

	go controlCenterInput(os.Stdin)
	go controlCenterCMD(sendCmds)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error: %s\n", err)
		}

		client := Client{
			conn: conn,
			ip: conn.RemoteAddr().String(),
			Status: true,
		}
		clients.Add(client)

		fmt.Printf("%s has connected\n", client.ip)

		wg.Add(1)
		go handleConn(client, &wg)

		go func() {
			wg.Wait()
			// close channels here
		}()
	}

}

func controlCenterInput(in io.Reader) {
	input := bufio.NewScanner(in)
	for input.Scan() {
		for _, client := range clients.Data {
			client.conn.Write(input.Bytes())
		}
	}
}

func controlCenterCMD(ch chan string) {

}

func handleConn(client Client, wg *sync.WaitGroup) {
	defer client.conn.Close()
	defer wg.Done()
	fmt.Println(client.ip)

	input := bufio.NewScanner(client.conn)
	for input.Scan() {
		if client.Status == true {
			fmt.Printf("From %s: %s\n", client.ip, input.Text())
		} else {
			client.SaveLogs(input.Text())
		}
	}

	fmt.Printf("%s has disconnected\n", client.ip)

	clients.Remove(client)
}
