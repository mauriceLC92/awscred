package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/mauriceLC92/awscred/credentials"
)

func CheckDefaultProfile() {
	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("an unexpected error occurred loading the default profile%v", err)
		return
	}

	// Create a new STS client
	svc := sts.NewFromConfig(cfg)

	// Call the "GetCallerIdentity" API to check if the credentials are valid
	_, err = svc.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		logError(err, "default")
	}
	logSuccess("default")
}

func CheckGivenProfile(profileName string) {
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

func logSuccess(profileName string) {
	greenTickEmoji := '\u2705'
	fmt.Printf("Profile \"%s\" - %c \n", profileName, greenTickEmoji)
}

func logError(err error, profileName string) {
	redCrossMarkEmoji := '\u274C'
	var oe *smithy.OperationError
	if errors.As(err, &oe) {
		errString := oe.Unwrap().Error()
		if strings.Contains(errString, "403") {
			fmt.Printf("Profile \"%s\" has invalid credentials - %c \n", profileName, redCrossMarkEmoji)
		} else {
			// The error is not related to the credentials
			fmt.Printf("An unexpected error occurred: %v", err)
		}
	}
}

func CheckCredentials(credentials []credentials.Credential) {
	CheckDefaultProfile()
	// check all beyond the default
	for _, cred := range credentials[1:] {
		CheckGivenProfile(cred.ProfileName)
	}
}
