package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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

	for _, cred := range credentials {
		fmt.Println("Profile - ", cred.profileName)
		fmt.Println("Access Key ID - ", cred.accessKeyId)
		fmt.Println("Secret Access Key - ", cred.accesskey)
	}
}
