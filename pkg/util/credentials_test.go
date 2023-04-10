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

func TestGetCredentialsProfiles(t *testing.T) {
	_, err := GetCredentialsProfiles()
	if err != nil {
		t.Fail()
	}
}

func TestGetConfigProfiles(t *testing.T) {
	val, err := GetConfigProfiles()
	if err != nil {
		t.Fail()
	}
	t.Log(val)
}

func TestGetAllProfiles(t *testing.T) {
	val, err := GetAllProfiles()
	if err != nil {
		t.Fail()
	}
	t.Log(val)
}

func TestClearProfile(t *testing.T) {
	err := ClearProfiles([]string{"yo"})
	if err != nil {
		t.Fail()
	}
}
