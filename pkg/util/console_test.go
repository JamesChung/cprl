package util_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/JamesChung/cprl/pkg/util"
)

func TestRead(t *testing.T) {
	v, _ := util.GetCredentials("default")
	fmt.Printf("%#v\n", v)
	js, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(string(js))
}

// func TestGetFederatedToken(t *testing.T) {
// 	u := url.URL{
// 		Scheme: "https",
// 		Host:   "google.com",
// 	}
// 	res, err := util.GetFederatedToken(u)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(res)
// }

func TestGenerateLoginURL(t *testing.T) {
	creds, err := util.GetCredentials("default")
	if err != nil {
		t.Fatal(err)
	}
	loginURL, err := util.GenerateLoginURL(creds, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(loginURL.String())
}
