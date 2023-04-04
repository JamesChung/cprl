package util_test

import (
	"fmt"
	"testing"

	"github.com/JamesChung/cprl/pkg/util"
)

func TestGetCredentialsFromCache(t *testing.T) {
	creds, err := util.GetCredentialsFromCache()
	if err != nil {
		t.Fail()
	}
	fmt.Println(creds)
}
