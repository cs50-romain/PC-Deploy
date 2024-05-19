package controlcenter

import (
	"fmt"
	"os"
	"text/template"
)

type Package struct {
	ComputerName	string
	TimeZone	string
	LocalAdmin	string
	LocalLolz	string
	Wifi		bool
	sSID		string
	sSIDLolz	string
}

func NewPackage(computerName, timeZone string) *Package {
	return &Package{
		ComputerName: computerName,
		TimeZone: timeZone,
		Wifi: false,
	}
}

func (p *Package) SetLocalAdmin(username, lolz string) {
	p.LocalAdmin = username
	p.LocalLolz = lolz
}

func (p *Package) SetWifi(ssid, lolz string) {
	p.sSID = ssid
	p.sSIDLolz = lolz
}

func (p *Package) WriteXML(buf []byte) error {
	file, err := os.OpenFile("/home/lettuce/go/PCDeploy/storage/clients/Autounattend.xml", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	t := template.Must(template.New("Autounattend.xml").Parse(string(buf)))	
	return t.Execute(file, p)
}

func (p *Package) ReadFile() ([]byte, error) {
	buf := make([]byte, 1000000000)
	file, err := os.Open("/home/lettuce/go/PCDeploy/deployment/Autounattend.xml")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	
	readBits, err := file.Read(buf)
	if err != nil {
		fmt.Printf("Bits read: %d", readBits)
		return nil, err
	}

	return buf, nil
}
