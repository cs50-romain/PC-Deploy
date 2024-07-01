package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var provisionCmd = &cobra.Command{
	Use: "provision",
	Short: "Provision a drive with the provision package for a client",
	Long: "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Provisionin")
	},
}

// Flags
var XmlFile string
var Client string
var Drive string
// POWER???????

func init() {
	cobra.OnFinalize(handleProvision)
	provisionCmd.Flags().StringVarP(&XmlFile, "xmlfile", "x", "", "autounnatend.xml file for provisioning PC")
	provisionCmd.Flags().StringVarP(&Client, "client", "c", "", "Name of the client for which we are provisioning.")
	provisionCmd.Flags().StringVarP(&Drive, "drive", "d", "", "Drive where the provisioning package will be copied to")
	provisionCmd.MarkFlagRequired("drive")
	provisionCmd.MarkFlagRequired("client")
	rootCmd.AddCommand(provisionCmd)
}

func handleProvision() {
	drive := provisionCmd.Flag("drive").Value.String()
	client := provisionCmd.Flag("client").Value.String()
	xmlFile := provisionCmd.Flag("xmlfile").Value.String()

	// check if client path exists
	if clientExists := exists("clients/" + client); clientExists != true {
		fmt.Fprintf(os.Stderr, "Client %s does not exist. Quitting!\n", client)
		return
	}

	// Copy client folder, xml file to the drive specified
	clientPath := "clients/" + client
	os.MkdirAll(drive + "/" + clientPath, 0750)
	copyDirs(clientPath, drive + "/" + clientPath, drive)

	copyFile(xmlFile, drive + "/autounattend.xml")
	//handlePowerOptions("nope")
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

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
