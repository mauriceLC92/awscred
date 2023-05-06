package awscred_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mauriceLC92/awscred"
)

// Parse reads a credentials file and returns a slice of Credentials.
func TestParseReadsCredentialsFileAndReturnSliceOfCredentials(t *testing.T) {
	t.Parallel()

	creds, err := awscred.Parse("testdata/test-credentials.txt")
	if err != nil {
		t.Fatal("error parsing file")
	}

	want := []awscred.Credential{
		{
			ProfileName: "test-1",
			AccessKeyId: "123",
			Accesskey:   "qwerty",
		},
		{
			ProfileName: "test-2",
			AccessKeyId: "321",
			Accesskey:   "qwerty-two",
		},
	}

	if !cmp.Equal(creds, want) {
		t.Error("file creds do not match expected credentials", cmp.Diff(want, creds))
	}
}

// Print reads a slice of credentials and prints the output to the terminal
func TestPrintTo_ReadsASliceOfCredentialsAndPrintsToGivenWriter(t *testing.T) {
	t.Parallel()
	buf := new(bytes.Buffer)
	creds := []awscred.Credential{
		{
			ProfileName: "test-1",
			AccessKeyId: "123",
			Accesskey:   "qwerty",
		},
	}
	awscred.PrintTo(buf, creds)
	want := "Profile - test-1\nAccess Key ID - 123\nSecret Access Key - qwerty\n"
	got := buf.String()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

// DeleteCredentialByProfile should error when given an invalid file path to credentials file
func TestDeleteCredentialByProfile_ErrorsWhenGivenInvalidPathToCredentialsFile(t *testing.T) {
	t.Parallel()

	credentialsFile := "testdata/some-random-path.txt"
	profileName := "something"
	err := awscred.DeleteCredentialByProfile(profileName, credentialsFile)
	if err == nil {
		t.Errorf("expected file %s not to exist but instead it was found", credentialsFile)
	}
}

// DeleteCredentialByProfile should delete a credential based on given profile name
func TestDeleteCredentialByProfile_DeletesACredentialByProfileName(t *testing.T) {
	t.Parallel()

	credentialsFile := "testdata/invalid-credentials.txt"
	profileName := "there"
	awscred.DeleteCredentialByProfile(profileName, credentialsFile)

	got, _ := os.ReadFile(credentialsFile)
	gotString := string(got)
	fmt.Printf("gotString: %v\n", gotString)

}

// DeleteConfigByProfile should error when given an invalid file path to credentials file
func TestDeleteConfigByProfile_ErrorsWhenGivenInvalidPathToConfigFile(t *testing.T) {
	t.Parallel()

	credentialsFile := "testdata/some-random-path.txt"
	profileName := "something"
	err := awscred.DeleteConfigByProfile(profileName, credentialsFile)
	if err == nil {
		t.Errorf("expected file %s not to exist but instead it was found", credentialsFile)
	}
}
