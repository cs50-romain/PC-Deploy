package main

import (
	"cs50-romain/pcdeploy/cmd/controlcenter"
	color "cs50-romain/pcdeploy/pkg"
	"fmt"
	"log"
)

func main() {
	// Start controlcenter & print intro
	fmt.Println(color.Fore(color.Yellow, "Welcome to PCDeploy. Let's make it easy to deploy PCs. Type help to view available commands!"))

	cc := controlcenter.ControlCenter{}

	if err := cc.Start(); err != nil {
		log.Fatal(err)
	}
}
