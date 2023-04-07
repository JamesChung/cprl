package client

import (
	"context"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	out, err := GetProfileConfig("default")
	if err != nil {
		t.Fail()
	}
	creds, err := out.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fail()
	}
	fmt.Println(creds)
}
