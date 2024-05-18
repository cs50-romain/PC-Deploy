package controlcenter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cs50-romain/tourdego/pkg/color"
	table "github.com/jedib0t/go-pretty/v6/table"
)

// I could actally create a client workspace. It would be a new shell, with new acceptable commands...// 
// Backing out of that would return us back to the initial shell
func SelectHandler(opts []string) (string, error) {
	if len(opts) == 0 {
		return "", fmt.Errorf(color.Bold + color.Red + "\tNot enough options given" + color.Reset)
	}
	
	if len(opts) > 1 {
		return "", fmt.Errorf(color.Bold + color.Red + "\tToo many options given" + color.Reset)
	}
	
	clientName := opts[0]
	if _, ok := clients[clientName]; !ok {
		return "", fmt.Errorf(color.Bold + color.Red + "\tClient does not exist. Use `show clients`to view available clients" + color.Reset)
	}

	fmt.Printf("\t%s%sYOU ARE NOW USING %s's WORKSPACE%s\n", color.Bold, color.Magenta, clientName, color.Reset)
	return clientName, nil
}

func ShowHandler(opts []string) {
	for _, opt := range opts {
		if opt != "clients" {
			fmt.Println(color.Bold + color.Red + "\tInvalid optional argument. Please use clients or packcages." + color.Reset)
			return
		}
	}
	// Show clients in a table of this format:
	// ID | Client Name | MORE LATER
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Client Name"})
	idx := 0
	for clientName := range clients {
		t.AppendSeparator()
		t.AppendRow([]interface{}{idx, clientName})
		idx++
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
		fmt.Println("\tPlease enter an option for this command; client or package?")
		return nil, fmt.Errorf("Invalid option")
	}

	for _, opt := range options {
		switch opt {
		case "client":
			// REFACTOR?
			// Ask questions and handle user input & create new client
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\tName of client: ")
			clientName, err := reader.ReadString('\n')
			clientName = strings.Trim(clientName, "\n")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("\tName of locations (main, second, third...): ")
			locations, err := reader.ReadString('\n')
			locations = strings.Trim(locations, "\n")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("\tEnter power options (format: _, _, _, _): ")
			powerOptions, err := reader.ReadString('\n')
			powerOptions = strings.Trim(powerOptions, "\n")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("\tEnter domain name: ")
			domainName, err := reader.ReadString('\n')
			domainName = strings.Trim(domainName, "\n")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			fmt.Print("\tEnter default PC name: ")
			pcname, err := reader.ReadString('\n')
			pcname = strings.Trim(pcname, "\n")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// Show how the client was created, with which config
			fmt.Println("\tYou have created a new client with the following configs:")
			fmt.Printf("\tName: %s Locations: %s Power Options: %s Domain Name: %s Default PC Name: %s\n", clientName, locations, powerOptions, domainName, pcname)

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
			fmt.Println("\tCreating package")
			return nil, nil
		default:
			createPackage()
		}
	}
	return nil, nil
}

func createPackage() {
	
}

func backHandler(opts []string) error {
	return nil
}
