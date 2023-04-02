package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	"gopkg.in/ini.v1"
)

type Credentials struct {
	AccessKeyID     string `json:"sessionId"`
	SecretAccessKey string `json:"sessionKey"`
	SessionToken    string `json:"sessionToken"`
}

func GetCredentials(profile string) (Credentials, error) {
	// Find user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Credentials{}, err
	}

	f, err := ini.Load(path.Join(homeDir, ".aws/credentials"))
	if err != nil {
		return Credentials{}, err
	}

	s := f.Section(profile)

	return Credentials{
		AccessKeyID:     s.Key("aws_access_key_id").String(),
		SecretAccessKey: s.Key("aws_secret_access_key").String(),
		SessionToken:    s.Key("aws_session_token").String(),
	}, nil
}

func StringifyCredentials(creds Credentials) (string, error) {
	b, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GenerateLoginURL(creds Credentials, isGov bool) (url.URL, error) {
	c, err := StringifyCredentials(creds)
	if err != nil {
		return url.URL{}, err
	}

	federationQuery := url.Values{}
	federationQuery.Add("Action", "getSigninToken")
	federationQuery.Add("DurationSeconds", "43200")
	federationQuery.Add("Session", c)
	federationURL := url.URL{
		Scheme:   "https",
		Host:     "signin.aws.amazon.com",
		Path:     "federation",
		RawQuery: federationQuery.Encode(),
	}

	token, err := GetFederatedToken(federationURL)
	if err != nil {
		return url.URL{}, err
	}

	loginQuery := url.Values{}
	loginQuery.Add("Action", "login")
	loginQuery.Add("Destination", "https://console.aws.amazon.com/")
	loginQuery.Add("SigninToken", token)
	loginQuery.Add("Issuer", "https://example.com")
	loginURL := url.URL{
		Scheme:   "https",
		Host:     "signin.aws.amazon.com",
		Path:     "federation",
		RawQuery: loginQuery.Encode(),
	}

	return loginURL, nil
}

func GetFederatedToken(u url.URL) (string, error) {
	c := http.Client{
		Timeout: time.Second * 5,
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
