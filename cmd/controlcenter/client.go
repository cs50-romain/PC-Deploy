package controlcenter

import (
	"encoding/json"
	"io"
	"strings"
)

type Client struct {
	Name		string
	Locations	[]string
	Automate	bool
	ConfigFile	*ClientConfig
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
		ConfigFile: configFile,
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
