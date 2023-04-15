package awscred_test

import (
	"bytes"
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
func TestPrintToReadsASliceOfCredentialsAndPrintsToGivenWriter(t *testing.T) {
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

// IsValidProfile takes a profile name and checks if valid or not
func TestIsValidProfileValidOrNot(t *testing.T) {
	// TODO - how does one test this?
}

// Delete Profile reads a credentials file and removes a profile if present
// func TestCleanRemovesInvalidProfiles(t *testing.T) {
// 	t.Parallel()

// 	awscred.DeleteProfile("test-2")

// }
