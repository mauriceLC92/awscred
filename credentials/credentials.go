package credentials

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Credential struct {
	ProfileName string
	AccessKeyId string
	Accesskey   string
}

func GetCredentials(file *os.File) []Credential {
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
		fmt.Fprintln(os.Stderr, "error reading from credentials file", err)
	}

	return credentials
}

func PrintCredentials(credentials []Credential) {
	for _, cred := range credentials {
		fmt.Println("Profile - ", cred.ProfileName)
		fmt.Println("Access Key ID - ", cred.AccessKeyId)
		fmt.Println("Secret Access Key - ", cred.Accesskey)
		fmt.Println("--------------------------------------------------------------")
	}
}
