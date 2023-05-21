package awscred

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
)

const (
	AWS_CREDENTIALS = "$HOME/.aws/credentials"
	AWS_CONFIG      = "$HOME/.aws/config"

	HELP_PRINT = "print - Display all the credentials found in '.aws/credentials`."
	HELP_CHECK = "check - Checks and displays if the credentials found in '.aws/credentials` are valid or not."
	HELP_CLEAN = "clean - Removes any invalid credentials found from both '.aws/credentials` and '.aws/config`."
	HELP_APPLY = "apply - Creates a shell with the given AWS profile exported so that any commands given run within the context of the desired profile."
	HELP_HELP  = "help - Display a menu of the commands available and their descriptions. \n"
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

func GenerateHelpMenu() string {
	return strings.Join([]string{HELP_PRINT, HELP_CHECK, HELP_CLEAN, HELP_HELP}, "\n")
}

func PrintHelpTo(w io.Writer, s string) {
	fmt.Fprint(w, s)
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
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		fmt.Printf("err: %v\n", err)
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
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	newData := []byte{}
	for scanner.Scan() {
		if scanner.Text() == fmt.Sprintf("[profile %s]", profileName) || scanner.Text() == fmt.Sprintf("[%s]", profileName) {
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

func Clean(credentials []Credential) {
	invalidProfiles := []string{}

	for _, cred := range credentials {
		err := IsValidProfile(cred.ProfileName)
		if err != nil {
			invalidProfiles = append(invalidProfiles, cred.ProfileName)
		}
	}

	for _, profile := range invalidProfiles {
		err := DeleteCredentialByProfile(profile, os.ExpandEnv(AWS_CREDENTIALS))
		if err != nil {
			fmt.Printf("err DeleteCredentialByProfile: %v\n", err)
		}
		err = DeleteConfigByProfile(profile, os.ExpandEnv(AWS_CONFIG))
		if err != nil {
			fmt.Printf("err DeleteConfigByProfile: %v\n", err)
		}
	}
}

func Apply(profileName string) {
	os.Setenv("AWS_PROFILE", profileName)
	// Be able to run `sls` (serverless cli) commands
	os.Setenv("SLS_INTERACTIVE_SETUP_ENABLE", "1")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userInput := scanner.Text()

		shellCommand := "sh"
		commandToRun := []string{"-c", userInput}

		command := exec.Command(shellCommand, commandToRun[0:]...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		err := command.Run()
		if err != nil {
			fmt.Printf("error running the command %s \n: %v", commandToRun[1], err)
		}
	}
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
