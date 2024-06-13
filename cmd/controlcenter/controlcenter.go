/*
This package will be used as the control center of the whole program. This is where the shell will runn, commands will be parse, data will be parse, added/removed/edited from "database".

This is the whole control flow of the program
*/
package controlcenter

import (
	"cs50-romain/pcdeploy/cmd/client"
	"cs50-romain/pcdeploy/cmd/controlcenter/workspace"
	"cs50-romain/pcdeploy/cmd/server"
	"fmt"
	"net"

	tourdego "github.com/cs50-romain/tourdego/pkg"
	"github.com/cs50-romain/tourdego/pkg/color"
)

var clients = make(map[string]*Client)

type ControlCenter struct {
	Commands	 map[string]func()
	Workspace	 bool
	clients		 map[string]*Client
	serv		 *server.Server
	currentWorkspaceName string // either <clientName> or <ip address>
	currentWorkspace *workspace.Workspace
}

func (c *ControlCenter) Start() error {
	// Add commands and start the shell
	prompt := "PCDeploy"
	sh := tourdego.NewShell(prompt + color.White + "> ")
	sh.SetPromptColor(color.Cyan)
	sh.AddCommand("show", &tourdego.Cmd{
		Name: "show",
		Help: "list things, usually whatever comes after list: list clients",
		Handler: func (s ...string) error {
			/*
			if c.InWorkspace() {
				// Check client commands and its handler
				return nil
			}
			*/
			c.ShowHandler(s)
			return nil
		},
	})

	sh.AddCommand("server", &tourdego.Cmd{
		Name: "server",
		Help: "Start the server",
		Handler: func (s ...string) error {
			if len(s) == 1 {
				if s[0] == "stop" {
					fmt.Println("Received command to stop server")
					c.serv.Stop()
				}
			} else {
				// CHECK IF SERVER ALREADY EXISTS 
				go func() {
					ser := server.NewServer("192.168.6.112", 6969)
					c.serv = ser
					c.serv.ClientComputers.Ips = make(map[string]*client.ClientComputer)
					err := ser.Run()
					if err != nil {
						fmt.Println(err)
					}
				}()
			}

			return nil
		},
	})

	// Stop command - stop server, stop conn

	sh.AddCommand("create", &tourdego.Cmd{
		Name: "create",
		Help: "create <option>; Create a client or package then answer questions.",
		Handler: func (s ...string) error {
			if c.InWorkspace() {
				// Check workspace commands and its handler
				c.currentWorkspace.HandleCommands("create", s...)
				return nil
			}
			// Read user input to find out what to create.
			client, err := CreateHandler(s)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if client != nil {
				clients[client.Name] = client
			}
			return nil
		},
	})
	
	sh.AddCommand("use", &tourdego.Cmd{
		Name: "use",
		Help: "Make a selection to do further work with that selection. - select <client #>",
		Handler: func (s ...string) error {
			if c.InWorkspace() {
				// Check client commands and its handler
				c.currentWorkspace.HandleCommands("create", s...)
				return nil
			}

			workspaceName, workspace, err := c.SelectHandler(s)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			c.currentWorkspaceName = workspaceName 
			c.currentWorkspace = workspace
			//c.clients[c.currentWorkspaceName].Workspace = workspace
			sh.SetPrompt(color.Cyan + prompt + color.White + "\\" + color.Yellow + workspaceName + color.White + "> ")
			return nil
		},
	})

	sh.AddCommand("back", &tourdego.Cmd{
		Name: "back",
		Help: "back; Back out of current selection",
		Handler: func (s ...string) error {
			if !c.InWorkspace() {
				fmt.Println("Not in a workspace ", c.Workspace)
				// Check client commands and its handler
				return nil
			}

			err := c.backHandler(s)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			sh.SetPromptColor(color.Cyan)
			sh.SetPrompt(prompt + color.White + "> " + color.Reset)
			sh.SetPromptColor(color.Cyan)
			return nil
		},
	})

	sh.AddCommand("kill", &tourdego.Cmd{
		Name: "kill",
		Help: "kill connection",
		Handler: func (s ...string) error {
			if c.InWorkspace() {
				// Check workspace commands and its handler
				c.currentWorkspace.HandleCommands("kill", s...)
				return nil
			}
			
			return nil
		},
	})

	// Init client array
	cls, err := InitClients()
	if err != nil {
		return err
	}
	clients = cls
	c.clients = cls

	//sh.RawMode = true // Couple things need fixed before being able to use raw mode

	if err := sh.Start(); err != nil {
		return err
	}

	return nil
}

func (c *ControlCenter) InWorkspace() bool {
	return c.Workspace
}

func getServerIP() string {
	conn, err := net.Dial("udp", "8.8.8.8")
	if err != nil {
		fmt.Println("Unable to get server ip")
		return ""
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
