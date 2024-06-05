package controlcenter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cs50-romain/tourdego/pkg/color"
	"cs50-romain/pcdeploy/cmd/controlcenter/workspace"
	table "github.com/jedib0t/go-pretty/v6/table"
)

// THIS DOES A LOT.... maybe break it down
func (c *ControlCenter) SelectHandler(opts []string) (string, *workspace.Workspace, error) {
	if len(opts) <= 1 {
		return "", nil, fmt.Errorf(color.Bold + color.Red + "\tNot enough options given" + color.Reset)
	}
	
	if len(opts) > 2 {
		return "", nil, fmt.Errorf(color.Bold + color.Red + "\tToo many options given" + color.Reset)
	}
	
	// Setup workspace
	option := opts[0]
	chosenWorkspace := opts[1]
	var returnedWorkspace *workspace.Workspace
	if option == "client" {
		if _, ok := clients[chosenWorkspace]; !ok {
			return "", nil, fmt.Errorf(color.Bold + color.Red + "\tClient does not exist. Use `show clients`to view available clients" + color.Reset)
		}
		client := clients[chosenWorkspace]

		returnedWorkspace = workspace.InitWorkspace(chosenWorkspace)
		client.Workspace = returnedWorkspace
		c.Workspace = true

		// Add commands related to client workspace
		client.Workspace.AddCommand("create", func(s ...string) {
			fmt.Println("Creating")
		})

		fmt.Printf("\t%s%sYOU ARE NOW USING %s's WORKSPACE%s\n", color.Bold, color.Magenta, chosenWorkspace, color.Reset)
	}

	if option == "conn" {
		fmt.Println(c.serv.ClientComputers.Ips)
		if _, ok := c.serv.ClientComputers.Ips[chosenWorkspace]; !ok {
			return "", nil, fmt.Errorf(color.Bold + color.Red + "\tConn does not exist. Use `show conns` to view active connections" + color.Reset)
		}
		conn := c.serv.ClientComputers.Ips[chosenWorkspace]

		returnedWorkspace = workspace.InitWorkspace(chosenWorkspace)
		conn.Workspace = returnedWorkspace
		c.Workspace = true

		// Add commands related to connection workspace
		conn.Workspace.AddCommand("listen", func(s ...string) {
			fmt.Println("Listening...")
		})

		conn.Workspace.AddCommand("kill", func(s ...string) {
			c.serv.ClientComputers.Remove(*conn)
		})

		fmt.Printf("\t%s%sYOU ARE NOW USING %s's WORKSPACE%s\n", color.Bold, color.Magenta, chosenWorkspace, color.Reset)
	}
	return chosenWorkspace, returnedWorkspace, nil
}

func (c *ControlCenter) ShowHandler(opts []string) {
	if len(opts) != 1 {
		fmt.Println(color.Bold + color.Red + "\tNot enough options given. Please input what to show (eg: clients, conns...)" + color.Reset)
		return
	}
	for _, opt := range opts {
		if opt != "clients" && opt != "conns" {
			fmt.Println(color.Bold + color.Red + "\tInvalid optional argument. Please use clients or packcages." + color.Reset)
			return
		}
	}

	opt := opts[0]
	if opt == "clients" {
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
	} else if opt == "conns" {
		// Show conns in a table of this format:
		// ID | IP | MORE LATER
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Client Name"})
		idx := 0
		for ip := range c.serv.ClientComputers.Ips {
			t.AppendSeparator()
			t.AppendRow([]interface{}{idx, ip})
			idx++
		}
		t.SetStyle(table.StyleBold)
		t.Render()	
	}
}

// create client:
// When creating a client, we need to create a default config json file.
// NEED TO FIND A WAY TO ACCEPT USER INPUT WHILE IN RAW MODE OF SHELL
func CreateHandler(options []string) (*Client, error){
	if len(options) == 0 {
		fmt.Println("\tPlease enter an option for this command; client or package?")
		return nil, fmt.Errorf("Invalid option")
	}

	for _, opt := range options {
		switch opt {
		case "client":
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

			fmt.Println("\tYou have created a new client with the following configs:")
			fmt.Printf("\tName: %s Locations: %s Power Options: %s Domain Name: %s Default PC Name: %s\n", clientName, locations, powerOptions, domainName, pcname)

			config := CreateConfig(powerOptions, pcname, domainName)
			client := CreateClient(clientName, locations, config)

			filepath := "/home/lettuce/go/PCDeploy/storage/clients/" + strings.Trim(client.Name, "\n") + ".json"
			file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			client.SaveFile(file)
			return client, nil
		case "package":
			// Will need to find out which client's workspace we are currently in.
			// So need to create workspaces.
			// For now default to Nomes

			// ADD USER INPUTS/QUESTIONS to create new package
			// MOST SHOULD BE OPTIONAL.

			// Display a progress bar here
			fmt.Println("\tCreating package")

			// Handle this goroutine better
			go func() (*Client, error) {
				err := createPackage()
			
				if err != nil {
					return nil, err
				}
				return nil, nil
			}()
			return nil, nil 
		default:
		}
	}
	return nil, nil
}

func createPackage() error {
	p := NewPackage("NOMES-2024", "Eastern Standard Time")
	p.SetLocalAdmin("MyAdmin", "IAmAdminHereBitches")
	buf, err := p.ReadFile()
	if err != nil {
		fmt.Println("Here is my error")
		return err
	}
	
	err = p.WriteXML(buf)
	if err != nil {
		fmt.Println("WriteXML error:")
		return err
	}
	return nil
}

func (c *ControlCenter) backHandler(opts []string) error {
	c.Workspace = false
	return nil
}
