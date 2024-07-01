package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "pcdeploy",
	Short: "PCDeploy is a pc deployment tool for Windows",
	Long: "Nah",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are running pcdeploy")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

/*
import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

const STRNIL = ""

type Flags struct {
	port		string
	xmlFile		string
	drive		string
	client		string
	power		string
	computerName	string
	username	string
	password	string
	timezone	string
}

const TEMPLATE_XMLFILE_PATH = "./deployment/template_autounattend.xml"

func main() {
	action := flag.String("action", "", "Which action to perform.")
	portPtr := flag.String("port", "6969", "Port for server listening")
	clientPtr := flag.String("client", "", "The work client folder, which includes the automate agent and any other exe")
	xmlFilePtr := flag.String("xmlfile", "", "autounattend.xml file for provisioning USB")
	powerPtr := flag.String("power", "0,0,0,0,0,0", "Power Settings. Format: 0,0,0,0,0,0")
	usernamePtr := flag.String("username", "admin", "Local admin username")
	passPtr := flag.String("pass", "admin", "Local admin password")
	computerNamePtr := flag.String("computerName", "PCNAME", "Computer name")
	timeZonePtr := flag.String("tz", "Standard Eastern Time", "Provide timezone for the computer")
	drivePtr := flag.String("drive", "", "Drive to provision")

	flag.Parse()

	flags := Flags{
		port: *portPtr,
		xmlFile: *xmlFilePtr,
		drive: *drivePtr,
		client: *clientPtr,
		power: *powerPtr,
		computerName: *computerNamePtr,
		username: *usernamePtr,
		password: *passPtr,
		timezone: *timeZonePtr,
	}
	
	if err := handleAction(*action, flags); err != nil {
		fmt.Println(err)
		return 
	}
}

func handleAction(action string, flags Flags) error {
	if action != "provision" && action != "createxml" && action != "server" {
		return fmt.Errorf("Action is not correct. Please choose either provision, create-xml, or server.")
	}

	if action == "provision" {
		return handleProvision(flags.client, flags.drive, flags.xmlFile, flags)
	}

	if action == "createxml" {
		return handleXml(flags)
	}

	if action == "server" {
		return handleServer()
	}

	return nil
}

func handleProvision(client string, drive string, xmlFile string, flags Flags) error {
	if client == "" {
		return fmt.Errorf("Client cannot be nil")
	}
	if drive == "" {
		return fmt.Errorf("Drive cannot be nil")
	}
	
	if xmlFile == "" {
		return fmt.Errorf("xmlfile cannot nil")
	}
	
	// Copy client folder, xml file to the drive specified
	clientPath := strings.TrimLeft(client, ".")
	os.MkdirAll(drive + clientPath, 0750)
	copyDirs(client, drive + clientPath, drive)

	copyFile(xmlFile, drive + "/autounattend.xml")
	handlePowerOptions(flags.power)

	return nil
}

func handleXml(flags Flags) error {
	if flags.drive == "" {
		return fmt.Errorf("Drive cannot be nil")
	}

	xmlPackage := NewPackage(flags.computerName, flags.timezone)
	xmlPackage.SetLocalAdmin(flags.username, flags.password)

	src, err := os.Open(TEMPLATE_XMLFILE_PATH)
	defer src.Close()
	if err != nil {
		return err
	}
	buf, err := xmlPackage.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return err
	}

	file, err := os.OpenFile(flags.drive + "autounattend.xml", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = xmlPackage.WriteXML(buf, file)
	if err != nil {
		return err
	}
	return nil
}

func handlePowerOptions(powerOptions string) error {
	powerOptionsSplit := strings.Split(powerOptions, ",")
	monitorTimeoutAC, monitorTimeoutDC := powerOptionsSplit[0], powerOptionsSplit[1]
	standbyTimeoutAC, standbyTimeoutDC := powerOptionsSplit[2], powerOptionsSplit[3]
	hibernateTimeoutAC, hibernateTimeoutDC := powerOptionsSplit[4], powerOptionsSplit[5]

	powerFile, err := os.ReadFile("./deployment/setpoweroption.bat")
	if err != nil {
		fmt.Println(err)
		return err
	}

	powers := []string{monitorTimeoutAC, monitorTimeoutDC, standbyTimeoutAC, standbyTimeoutDC, hibernateTimeoutAC, hibernateTimeoutDC}

	lines := strings.Split(string(powerFile), "\n")
	idx := 0
	for lineIdx, line := range lines {
		if len(line) == 0 {
			continue
		}
		if string(line[0]) != "p" {
			continue
		}
		lineSplit := strings.Split(line, " ")
		lineSplit[len(lineSplit)-1] = powers[idx]
		line = strings.Join(lineSplit, " ")
		lines[lineIdx] = line
		idx++
	}

	output := strings.Join(lines, "\n")
	fmt.Println(output)
	err = os.WriteFile("./deployment/setpoweroption.bat", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func handleServer(flags ...string) error {
	fmt.Println("Server not available yet")
	return nil
}

func copyDirs(sourceDirPath, destDirPath, drivePath string) error {
	entries, err := os.ReadDir(sourceDirPath)
	if err != nil {
		fmt.Println(err)
		return err 
	}

	for _, entry := range entries {
		entryPath := sourceDirPath + "/" + entry.Name()
		if entry.IsDir() {
			os.MkdirAll(destDirPath + "/" + entry.Name(), 0750)
			go copyDirs(entryPath, destDirPath + "/" + entry.Name(), drivePath)
		} else {
			if entry.Name() == "autounattend.xml" {
				copyFile(entryPath, drivePath + "/" + entry.Name())
			} else {
				copyFile(entryPath, destDirPath + "/" + entry.Name())
			}
		}
	}
	return nil
}

func copyFile(srcFile, destFile string) error {
	osrcFile, err := os.Open(srcFile)
	if err != nil {
		fmt.Println("Error opening source file: ", err)
		return err
	}
	defer osrcFile.Close()
	odstFile, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error with destination file: ", err)
		return err
	}
	defer odstFile.Close()
	return Copy(osrcFile, odstFile)
}

func Copy(source io.Reader, dest io.Writer) error {
	_, err := io.Copy(dest, source)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

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

func (p *Package) WriteXML(buf []byte, dest io.Writer) error {
	t := template.Must(template.New("autounattend.xml").Parse(string(buf)))	
	return t.Execute(dest, p)
}

func (p *Package) ReadFile(src io.Reader) ([]byte, error) {
	buf, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
*/
