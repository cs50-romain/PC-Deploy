package controlcenter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Client struct {
	Name		string `json:"Name"`
	Locations	[]string `json:"Locations"`
	Automate	bool `json:"Automate"`
	ConfigFile	ClientConfig
}

type ClientConfig struct {
	OOBE		bool
	PowerOption	[]string
	PCName		string
	DomainName	string
	GLOBs		[]string	
	RemoveBloatware bool
}

func CreateClient(name, locations string, configFile *ClientConfig) *Client {
	return &Client{
		Name: name,
		Locations: strings.Split(locations, ", "),
		Automate: true,
		ConfigFile: *configFile,
	}
}

func CreateConfig(powerOpt, pcname, domainname string) *ClientConfig {
	return &ClientConfig{
		OOBE: true,
		PowerOption: strings.Split(powerOpt, ", "), // Needs better parsing.
		PCName: pcname,
		DomainName: domainname,
		GLOBs: []string{"Google Chrome", "Adobe Reader DC", "Firefox"}, // Should ask next time
		RemoveBloatware: true,
	}
}

func (c *Client) SaveFile(w io.Writer) error {
	//path := "/home/lettuce/go/PCDeploy/storage/clients/"
	fileContents, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = w.Write(fileContents)
	//encoder := json.NewEncoder(w)
	//err = encoder.Encode(fileContents)
	if err != nil {
		return err
	}

	return nil
}

func getClientFromFile(filename string) (*Client, error) {
	// Open file
	filepath := "/home/lettuce/go/PCDeploy/storage/clients/" + filename
	fmt.Println(filename, filepath)
	byt, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var client Client
	if err = json.Unmarshal(byt, &client); err != nil {
		return nil, err
	}

	return &client, nil
}

func InitClients() (map[string]Client, error) {
	var cls = make(map[string]Client)

	// Read directory where all clients.json 
	dirs, err := os.ReadDir("/home/lettuce/go/PCDeploy/storage/clients/")
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		client, err := getClientFromFile(dir.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		cls[client.Name] = *client
	}
	return cls, nil
}
