package credentials_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mauriceLC92/go-awscred/credentials"
)

// Parse reads a credentials file and returns a slice of Credentials.
func TestParseReadsCredentialsFileAndReturnSliceOfCredentials(t *testing.T) {
	t.Parallel()

	creds, err := credentials.Parse("testdata/test-credentials.txt")
	if err != nil {
		t.Fatal("error parsing file")
	}

	want := []credentials.Credential{
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

// Refactor to call the parse function instead.
