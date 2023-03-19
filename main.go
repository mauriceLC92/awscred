package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauriceLC92/go-awscred/aws"
	"github.com/mauriceLC92/go-awscred/credentials"
)

func main() {
	file, err := os.Open(os.ExpandEnv("$HOME/.aws/credentials"))
	if err != nil {
		log.Fatalln("error opening credentials file:", err)
		return
	}
	defer file.Close()
	creds := credentials.GetCredentials(file)

	if len(os.Args) > 1 {
		commandLineArg := os.Args[1]
		fmt.Printf("commandLineArg: %v\n", commandLineArg)

		switch strings.ToLower(commandLineArg) {
		case "print":
			credentials.PrintCredentials(creds)
		case "check":
			aws.CheckCredentials(creds)
		case "apply":
			fmt.Println("applying profile...")
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
