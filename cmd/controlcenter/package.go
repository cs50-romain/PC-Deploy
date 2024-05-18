package controlcenter

import "encoding/xml"

//"encoding/xml"

// CREATE FUNCTION NAMES AND STRUCTS BEFORE ANYTHING
type Package struct {
	ComputerName	string
	timeZone	string
	localAdmin	string
	localLolz	string
	Wifi		bool
	sSID		string
	sSIDLolz	string
}

func NewPackage(computerName, timeZone string) *Package {
	return &Package{
		ComputerName: computerName,
		timeZone: timeZone,
	}
}

func (p *Package) SetLocalAdmin(username, lolz string) {
	p.localAdmin = username
	p.localLolz = lolz
}

func (p *Package) SetWifi(ssid, lolz string) {
	p.sSID = ssid
	p.sSIDLolz = lolz
}

func (p *Package) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return nil
}

func (p *Package) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}
