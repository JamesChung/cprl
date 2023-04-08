package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestWriteCredentials(t *testing.T) {
	err := WriteCredentials("test", aws.Credentials{
		AccessKeyID:     "AccessKeyID2",
		SecretAccessKey: "SecretAccessKey2",
		SessionToken:    "SessionToken2",
	})
	if err != nil {
		t.Fail()
	}
}

func TestGetProfiles(t *testing.T) {
	_, err := GetProfiles()
	if err != nil {
		t.Fail()
	}
}

func TestClearProfile(t *testing.T) {
	err := ClearProfile("yo")
	if err != nil {
		t.Fail()
	}
}
