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

const (
	AWS_CREDENTIALS = "$HOME/.aws/credentials"
	AWS_CONFIG      = "$HOME/.aws/config"
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

func IsValidProfile(profileName string) (err error) {
	// When struggling to think what to test, break down the function
	// into chunks of what the function is doing and what is essential to it performing it's function.
	// Imagine what bugs there would be and then write tests to detect them.
	// Example is if the svc is nil, we can't do anything with that!
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
		return err
	}
	return nil
}

func skipCredentials(s *bufio.Scanner) {
	s.Scan()
	s.Scan()
}

func skipConfig(s *bufio.Scanner) {
	s.Scan()
}

func DeleteCredentialByProfile(profileName, filePath string) error {
	// not sure how to know the FileMode for this?
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	newData := []byte{}
	for scanner.Scan() {
		if scanner.Text() == fmt.Sprintf("[%s]", profileName) {
			// skip the next two lines after finding the profile names
			skipCredentials(scanner)
			continue
		}
		newData = append(newData, append(scanner.Bytes(), "\n"...)...)
	}

	err = os.WriteFile(filePath, newData, 0777)
	if err != nil {
		return err
	}
	return nil
}

func DeleteConfigByProfile(profileName, filePath string) error {
	// not sure how to know the FileMode for this?
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	newData := []byte{}
	for scanner.Scan() {
		if scanner.Text() == fmt.Sprintf("[%s]", profileName) {
			// skip the next two lines after finding the profile names
			skipConfig(scanner)
			continue
		}
		newData = append(newData, append(scanner.Bytes(), "\n"...)...)
	}

	err = os.WriteFile(filePath, newData, 0777)
	if err != nil {
		return err
	}
	return nil
}

func CheckCredentials(credentials []Credential) {
	for _, cred := range credentials {
		err := IsValidProfile(cred.ProfileName)
		if err != nil {
			logError(err, cred.ProfileName)
			continue
		}
		logSuccess(cred.ProfileName)
	}
}

func ApplyProfile(profileName string) {
	cmd := exec.Command("export", fmt.Sprintf("AWS_PROFILE=%s", profileName))
	err := cmd.Run()
	if err != nil {
		logError(err, profileName)
	}
}

func Clean(credentials []Credential) {
	invalidProfiles := []string{}

	for _, cred := range credentials {
		err := IsValidProfile(cred.ProfileName)
		if err != nil {
			invalidProfiles = append(invalidProfiles, cred.ProfileName)
		}
	}

	for _, profile := range invalidProfiles {
		DeleteCredentialByProfile(profile, AWS_CREDENTIALS)
		DeleteConfigByProfile(profile, AWS_CONFIG)
	}
}

func Apply(profileName string, command []string) {
	comm := command[0]
	os.Setenv("AWS_PROFILE", profileName)
	output := strings.Join([]string{"AWS_PROFILE", "=", os.ExpandEnv("$AWS_PROFILE")}, "")
	cmd := exec.Command(comm, output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logError(err, profileName)
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
