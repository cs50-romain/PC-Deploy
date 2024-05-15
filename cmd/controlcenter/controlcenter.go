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

type ControlCenter struct {

}

func (c *ControlCenter) Start() error {
	// Add commands and start the shell
	sh := tourdego.NewShell("PCDeploy >")
	sh.SetPromptColor(color.Cyan)
	sh.AddCommand("show", &tourdego.Cmd{
		Name: "show",
		Help: "list things, usually whatever comes after list: list clients",
		Handler: func (...string) error {
			fmt.Println("Listing options")
			return nil
		},
	})

	sh.AddCommand("create", &tourdego.Cmd{
		Name: "create",
		Help: "create <option>; Create a client or package then answer questions.",
		Handler: func (...string) error {
			// Read user input to find out what to create.
			fmt.Println("Creating what you have decided sir")
			return nil
		},
	})

	sh.RawMode = true

	if err := sh.Start(); err != nil {
		return err
	}

	return nil
}
