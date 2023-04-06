package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"
)

func AWSHostname(isGov bool) string {
	if isGov {
		return "amazonaws-us-gov.com"
	}
	return "aws.amazon.com"
}

func GenerateLoginURL(creds Credentials, isGov bool) (url.URL, error) {
	c, err := StringifyCredentials(creds)
	if err != nil {
		return url.URL{}, err
	}

	signinURL := fmt.Sprintf("signin.%s", AWSHostname(isGov))

	federationQuery := url.Values{}
	federationQuery.Add("Action", "getSigninToken")
	federationQuery.Add("DurationSeconds", "43200")
	federationQuery.Add("Session", c)
	federationURL := url.URL{
		Scheme:   "https",
		Host:     signinURL,
		Path:     "federation",
		RawQuery: federationQuery.Encode(),
	}

	token, err := GetFederatedToken(federationURL)
	if err != nil {
		return url.URL{}, err
	}

	loginQuery := url.Values{}
	loginQuery.Add("Action", "login")
	loginQuery.Add("Destination", fmt.Sprintf("https://console.%s/", AWSHostname(isGov)))
	loginQuery.Add("SigninToken", token)
	federationURL.RawQuery = loginQuery.Encode()

	return federationURL, nil
}

func GetFederatedToken(u url.URL) (string, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}

	res, err := c.Get(u.String())
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get federated token, status code %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	data := map[string]string{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return "", err
	}

	return data["SigninToken"], nil
}

func OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform: [%s]", runtime.GOOS)
	}
	if err != nil {
		return err
	}

	return nil
}
