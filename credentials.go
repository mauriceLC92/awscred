package awscred

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
)

type Credential struct {
	ProfileName string
	AccessKeyId string
	Accesskey   string
}

func (c Credential) String() string {
	return fmt.Sprintf("Profile - %s\nAccess Key ID - %s\nSecret Access Key - %s\n", c.ProfileName, c.AccessKeyId, c.Accesskey)
}

func PrintTo(w io.Writer, credentials []Credential) {
	for _, cred := range credentials {
		fmt.Fprint(w, cred.String())
		fmt.Println("--------------------------------------------------------------")
	}
}

func Parse(filePath string) ([]Credential, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentProfile := Credential{}

	var credentials []Credential
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "[") {
			profileName := strings.Trim(line, "[]")
			currentProfile.ProfileName = profileName
			continue
		}
		if strings.HasPrefix(line, "aws_access_key_id") {
			currentProfile.AccessKeyId = strings.TrimPrefix(line, "aws_access_key_id = ")
			continue
		}

		if strings.HasPrefix(line, "aws_secret_access_key") {
			currentProfile.Accesskey = strings.TrimPrefix(line, "aws_secret_access_key = ")
			credentials = append(credentials, currentProfile)
			currentProfile = Credential{}
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return credentials, nil
}

func IsValidProfile(profileName string) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile(profileName),
	)
	if err != nil {
		fmt.Printf("An unexpected error occurred loading profile %s: %v", profileName, err)
		return
	}

	// Create a new STS client
	svc := sts.NewFromConfig(cfg)

	// Call the "GetCallerIdentity" API to check if the credentials are valid
	_, err = svc.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		logError(err, profileName)
		return
	}
	logSuccess(profileName)
}

func CheckCredentials(credentials []Credential) {
	for _, cred := range credentials {
		IsValidProfile(cred.ProfileName)
	}
}

func ApplyProfile(profileName string) {
	cmd := exec.Command("export", fmt.Sprintf("AWS_PROFILE=%s", profileName))
	err := cmd.Run()
	if err != nil {
		logError(err, profileName)
	}
}

// What I initially tried
func Apply(filePath string) {
	cmd := exec.Command("source", filePath)
	err := cmd.Run()
	if err != nil {
		logError(err, filePath)
	}
}

// This won't work either since this function will run within a child process
func ApplyIt(filePath string) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("source %s", filePath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logError(err, filePath)
	}
}

// A Go program is unable to update the environment variable of the shell that invoked it. The Go program
// is run in a child process of the shell. Child process cannot modify the environment of it's parent.
// The below generates a script which can then be sourced by the user running that within the shell.

func GenerateProfileScript(profileName string) {
	script := fmt.Sprintf(`#!/bin/bash
	export AWS_PROFILE=%s
	`, profileName)

	scriptFilePath := filepath.Join(os.TempDir(), "set_aws_profile.sh")
	err := ioutil.WriteFile(scriptFilePath, []byte(script), 0755)
	if err != nil {
		log.Fatal("Error writing script file:", err)
	}

	fmt.Printf("Script created: %s\n", scriptFilePath)
	fmt.Println("To update the environment variable, run:")
	fmt.Printf("source %s\n", scriptFilePath)
	ApplyIt(scriptFilePath)
}

func logSuccess(profileName string) {
	greenTickEmoji := '\u2705'
	fmt.Printf("Profile \"%s\" %c \n", profileName, greenTickEmoji)
}

func logError(err error, profileName string) {
	redCrossMarkEmoji := '\u274C'
	var oe *smithy.OperationError
	if errors.As(err, &oe) {
		errString := oe.Unwrap().Error()
		if strings.Contains(errString, "403") {
			fmt.Printf("Profile \"%s\" has invalid credentials %c \n", profileName, redCrossMarkEmoji)
		} else {
			// The error is not related to the credentials
			fmt.Printf("An unexpected error occurred: %v", err)
		}
	}
}
