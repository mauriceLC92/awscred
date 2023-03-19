package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauriceLC92/go-awscred/aws"
)

type credential struct {
	profileName string
	accessKeyId string
	accesskey   string
}

func main() {
	file, err := os.Open(os.ExpandEnv("$HOME/.aws/credentials"))
	if err != nil {
		log.Fatalln("error opening credentials file:", err)
		return
	}
	defer file.Close()
	credentials := getCredentials(file)

	if len(os.Args) > 1 {
		commandLineArg := os.Args[1]
		fmt.Printf("commandLineArg: %v\n", commandLineArg)

		switch strings.ToLower(commandLineArg) {
		case "print":
			printCredentials(credentials)
		case "check":
			aws.CheckDefaultProfile()
			for _, cred := range credentials[1:] {
				aws.CheckGivenProfile(cred.profileName)
			}
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

func getCredentials(file *os.File) []credential {
	scanner := bufio.NewScanner(file)
	currentProfile := credential{}

	var credentials []credential
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "[") {
			profileName := strings.Trim(line, "[]")
			currentProfile.profileName = profileName
			continue
		}
		if strings.HasPrefix(line, "aws_access_key_id") {
			currentProfile.accessKeyId = strings.TrimPrefix(line, "aws_access_key_id = ")
			continue
		}

		if strings.HasPrefix(line, "aws_secret_access_key") {
			currentProfile.accesskey = strings.TrimPrefix(line, "aws_secret_access_key = ")
			credentials = append(credentials, currentProfile)
			currentProfile = credential{}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from credentials file", err)
	}

	return credentials
}

func printCredentials(credentials []credential) {
	for _, cred := range credentials {
		fmt.Println("Profile - ", cred.profileName)
		fmt.Println("Access Key ID - ", cred.accessKeyId)
		fmt.Println("Secret Access Key - ", cred.accesskey)
		fmt.Println("--------------------------------------------------------------")
	}
}
