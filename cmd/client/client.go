package client

import (
	"cs50-romain/pcdeploy/cmd/controlcenter/workspace"
	"net"
	"regexp"
)

type ClientComputer struct {
	Conn 	net.Conn
	Ip 	string
	Status  bool		// foreground = true.
	Logs	[]string
	Workspace *workspace.Workspace
	Rcv	<-chan string
}

type ClientComputers struct {
	Ips   map[string]*ClientComputer
	Conns []*ClientComputer
}

func NewClientComputer(conn net.Conn, ip string) *ClientComputer {
	return &ClientComputer{
		Conn: conn,
		Ip: ip,
		Status: true,
		Logs: []string{},
		Rcv: make(<-chan string),
	}
}

func (c *ClientComputer) SaveLogs(log string) {
	c.Logs = append(c.Logs, log)
}

func (c *ClientComputers) Remove(client ClientComputer) {
	// Remove client from the array
	for i, conns := range c.Conns {
		if client.Ip == conns.Ip {
			c.Conns = append(c.Conns[:i], c.Conns[i+1:]...)	
		}
	}
}

func (c *ClientComputers) Add(client *ClientComputer) {
	c.Ips[client.Ip] = client
	c.Conns = append(c.Conns, client)
}

func (c *ClientComputers) AllToBackground() {
	for _, cl := range c.Conns {
		cl.Status = false
	}
}

func (c *ClientComputer) ToForeground() {
	c.Status = true
}

func (c *ClientComputer) ToBackgrounnd() {
	c.Status = false
}

func isConnection(input string) bool {
	re := regexp.MustCompile("[0-9]+")
	if len(re.FindAllString(input, -1)) <= 0 {
		return false
	}
	return true
}
