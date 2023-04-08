package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauriceLC92/awscred"
)

func main() {
	creds, err := awscred.Parse(os.ExpandEnv("$HOME/.aws/credentials"))
	if err != nil {
		log.Fatalln("error parsing credentials file:", err)
	}

	if len(os.Args) > 1 {
		commandLineArg := os.Args[1]
		fmt.Printf("commandLineArg: %v\n", commandLineArg)

		switch strings.ToLower(commandLineArg) {
		case "print":
			awscred.PrintTo(os.Stdout, creds)
		case "check":
			awscred.CheckCredentials(creds)
		case "apply":
			profileName := os.Args[2]
			awscred.ApplyProfile(profileName)
		case "clean":
			fmt.Println("cleaning AWS profiles")
		case "help":
			fmt.Println("Displays a friendly message of the options available.")
		default:
			fmt.Println("Command not recognised. Please use 'help' to see the commands available to you.")
		}
	} else {
		fmt.Println("no commands given")
	}
}
