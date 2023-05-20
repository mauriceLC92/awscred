package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauriceLC92/awscred"
)

func main() {
	if len(os.Args) > 1 {
		commandLineArg := os.Args[1]

		switch strings.ToLower(commandLineArg) {
		case "print":
			creds := getCreds()
			awscred.PrintTo(os.Stdout, creds)
		case "check":
			creds := getCreds()
			awscred.CheckCredentials(creds)
		case "apply":
			profileName := os.Args[2]
			awscred.Apply(profileName)
		case "clean":
			creds := getCreds()
			awscred.Clean(creds)
		case "help":
			log.Println("Displays a friendly message of the options available.")
		default:
			fmt.Println("Command not recognised. Please use 'help' to see the commands available to you.")
		}
	} else {
		log.Println("no commands given")
	}
}

func getCreds() []awscred.Credential {
	creds, err := awscred.Parse(os.ExpandEnv(awscred.AWS_CREDENTIALS))
	if err != nil {
		log.Fatalln("error parsing credentials file:", err)
	}
	return creds
}
