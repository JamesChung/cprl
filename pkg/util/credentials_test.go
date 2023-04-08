package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func Test(t *testing.T) {
	err := WriteCredentials("test", aws.Credentials{
		AccessKeyID:     "AccessKeyID2",
		SecretAccessKey: "SecretAccessKey2",
		SessionToken:    "SessionToken2",
	})
	if err != nil {
		t.Fail()
	}
}
