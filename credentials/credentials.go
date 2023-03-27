package credentials

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Credential struct {
	ProfileName string
	AccessKeyId string
	Accesskey   string
}

func PrintTo(w io.Writer, credentials []Credential) {
	for _, cred := range credentials {
		fmt.Fprintln(w, fmt.Sprint("Profile - ", cred.ProfileName))
		fmt.Fprintln(w, fmt.Sprint("Access Key ID - ", cred.AccessKeyId))
		fmt.Fprintln(w, fmt.Sprint("Secret Access Key - ", cred.Accesskey))
		// fmt.Println("--------------------------------------------------------------")
	}
}

func PrintCredentials(credentials []Credential) {
	for _, cred := range credentials {
		fmt.Println("Profile - ", cred.ProfileName)
		fmt.Println("Access Key ID - ", cred.AccessKeyId)
		fmt.Println("Secret Access Key - ", cred.Accesskey)
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
