package client

import "net"

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

func (c *Clients) AllToBackground() {
	for _, cl := range c.Data {
		cl.Status = false
	}
}

