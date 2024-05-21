/*
This package will be used as the control center of the whole program. This is where the shell will runn, commands will be parse, data will be parse, added/removed/edited from "database".

This is the whole control flow of the program
*/
package controlcenter

import (
	"fmt"

	tourdego "github.com/cs50-romain/tourdego/pkg"
	"github.com/cs50-romain/tourdego/pkg/color"
)

var clients = make(map[string]Client)

type ControlCenter struct {
	Commands map[string]func()
	Workspace	bool
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
			if c.InWorkspace() {
				// Check client commands and its handler
				return nil
			}
			ShowHandler(s)
			return nil
		},
	})

	sh.AddCommand("create", &tourdego.Cmd{
		Name: "create",
		Help: "create <option>; Create a client or package then answer questions.",
		Handler: func (s ...string) error {
			if c.InWorkspace() {
				// Check client commands and its handler
				return nil
			}
			// Read user input to find out what to create.
			client, err := CreateHandler(s)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if client != nil {
				clients[client.Name] = *client
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
				return nil
			}

			clientName, err := c.SelectHandler(s)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			sh.SetPrompt(color.Cyan + prompt + color.White + "\\" + color.Yellow + clientName + color.White + "> ")
			return nil
		},
	})

	sh.AddCommand("back", &tourdego.Cmd{
		Name: "back",
		Help: "back; Back out of current selection",
		Handler: func (s ...string) error {
			if c.InWorkspace() {
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

	// Init client array
	cls, err := InitClients()
	if err != nil {
		return err
	}
	clients = cls

	//sh.RawMode = true // Couple things need fixed before being able to use raw mode

	if err := sh.Start(); err != nil {
		return err
	}

	return nil
}

func (c *ControlCenter) InWorkspace() bool {
	return c.Workspace
}
