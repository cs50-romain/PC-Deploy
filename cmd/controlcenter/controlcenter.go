/*
This package will be used as the control center of the whole program. This is where the shell will runn, commands will be parse, data will be parse, added/removed/edited from "database".

This is the whole control flow of the program
*/
package controlcenter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tourdego "github.com/cs50-romain/tourdego/pkg"
	"github.com/cs50-romain/tourdego/pkg/color"
	table "github.com/jedib0t/go-pretty/v6/table"
)

var clients []Client

type ControlCenter struct {
	Commands map[string]func()
}

func (c *ControlCenter) Start() error {
	// Add commands and start the shell
	sh := tourdego.NewShell("PCDeploy" + color.White + "> ")
	sh.SetPromptColor(color.Cyan)
	sh.AddCommand("show", &tourdego.Cmd{
		Name: "show",
		Help: "list things, usually whatever comes after list: list clients",
		Handler: func (s ...string) error {
			ShowHandler(s)
			return nil
		},
	})

	sh.AddCommand("create", &tourdego.Cmd{
		Name: "create",
		Help: "create <option>; Create a client or package then answer questions.",
		Handler: func (s ...string) error {
			// Read user input to find out what to create.
			client, err := CreateHandler(s)
			if err != nil {
				return err
			}
			clients = append(clients, *client)
			return nil
		},
	})

	//sh.RawMode = true // Couple things need fixed before being able to use raw mode

	if err := sh.Start(); err != nil {
		return err
	}

	return nil
}

func ShowHandler(opts []string) {
	for _, opt := range opts {
		if opt != "clients" {
			fmt.Println("Invalid optional argument. Please use clients or packcages.")
			return
		}
	}
	// Show clients in a table of this format:
	// ID | Client Name | MORE LATER
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Client Name"})
	for idx, client := range clients {
		t.AppendSeparator()
		t.AppendRow([]interface{}{idx, client.Name})
	}
	t.SetStyle(table.StyleBold)
	t.Render()
}

// create client:
// When creating a client, we need to create a default config json file.
// User will need to answer questions.
// NEED TO FIND A WAY TO ACCEPT USER INPUT WHILE IN RAW MODE OF SHELL
func CreateHandler(options []string) (*Client, error){
	if len(options) == 0 {
		fmt.Println("Please enter an option for this command; client or package?")
		return nil, nil
	}

	for _, opt := range options {
		switch opt {
		case "client":
			// Ask questions and handle user input & create new client
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("   Name of client: ")
			clientName, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("   Name of locations (main, second, third...): ")
			locations, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("   Enter power options (format: _, _, _, _): ")
			powerOptions, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("   Enter domain name: ")
			domainName, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("   Enter default PC name: ")
			pcname, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// SHow how the client was created, with which config
			fmt.Println("   You have created a new client with the following configs:")
			fmt.Printf("   Name: %s Locations: %s Power Options: %s Domain Name: %s Default PC Name: %s", clientName, locations, powerOptions, domainName, pcname)

			config := CreateConfig(powerOptions, pcname, domainName)
			client := CreateClient(clientName, locations, config)

			// Add client to array

			// Save client to json file
			filepath := "/home/lettuce/go/PCDeploy/storage/clients/" + strings.Trim(client.Name, "\n") + ".json"
			file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			client.SaveFile(file)
			return client, nil
		case "package":
			fmt.Println("   Creating package")
			return nil, nil
		default:
			fmt.Println("Cannot create for this option. Please create a client or package.")
		}
	}
	return nil, nil
}
