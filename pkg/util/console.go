package util

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"gopkg.in/ini.v1"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
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

// func GenerateLoginURL(creds Credentials, isGov bool) (url.URL, error) {
// 	c, err := json.Marshal(creds)
// 	if err != nil {
// 		return url.URL{}, err
// 	}
// 	federationURL := url.URL{
// 		Scheme: "https",
// 		Host:   "signin.aws.amazon.com",
// 		Path:   "federation",
// 		RawQuery: fmt.Sprintf(
// 			"Action=getSigninToken&DurationSeconds=43200&Session=%s",
// 			string(c),
// 		),
// 	}

// 	token, err := GetFederatedToken(federationURL)
// 	if err != nil {
// 		return url.URL{}, nil
// 	}

// 	loginURL := url.URL{
// 		Scheme: "https",
// 		Host:   "",
// 	}
// 	return federationURL, nil
// }

func GetFederatedToken(u url.URL) (string, error) {
	c := http.Client{
		Timeout: time.Second * 3,
	}
	res, err := c.Get(u.String())
	if err != nil {
		return "", err
	}
	token, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
