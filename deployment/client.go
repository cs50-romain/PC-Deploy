package deployment

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func main(){
	Run()
}

func Run() {
	conn, err := net.Dial("tcp", "192.168.6.112:6969")
	if err != nil {
		panic(err)
	}

	client := &Client{conn}
	client.write("Hello there server\n")
	runPowershellScripts(client)
	client.write("Goodbye!\n")
}

func (c *Client) write(msg string) {
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}

func runPowershellScripts(client *Client) {
	// Run the runPowershellScripts\
	client.write("Ran script 1\n")
	time.Sleep(5 * time.Second)
	client.write("Ran script 2\n")
	time.Sleep(5 * time.Second)
	client.write("Ran script 3\n")
	time.Sleep(5 * time.Second)
	client.write("Ran script 4\n")
	time.Sleep(5 * time.Second)
	client.write("Ran script 5\n")
	time.Sleep(5 * time.Second)
	client.write("Ran script 6\n")
}
