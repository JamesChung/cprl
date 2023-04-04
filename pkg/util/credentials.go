package util

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"gopkg.in/ini.v1"
)

type Credentials struct {
	AccessKeyID     string `json:"sessionId"`
	SecretAccessKey string `json:"sessionKey"`
	SessionToken    string `json:"sessionToken"`
}

type CacheCredentials struct {
	ProviderType string
	Credentials  struct {
		AccessKeyId     string
		SecretAccessKey string
		SessionToken    string
		Expiration      string
	}
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

// func GetCacheCredentials(profile string) (Credentials, error) {

// }

func InitCredentialsCache(profile string) {
	//
}

func GetCredentialsFromCache() ([]Credentials, error) {
	creds := make([]Credentials, 0, 10)
	// Find user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cacheDir := path.Join(homeDir, ".aws/cli/cache")
	dir, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}
	for _, d := range dir {
		f, err := os.ReadFile(path.Join(cacheDir, d.Name()))
		if err != nil {
			return nil, err
		}
		cc := CacheCredentials{}
		err = json.Unmarshal(f, &cc)
		if err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, cc.Credentials.Expiration)
		if err != nil {
			return nil, err
		}
		if t.After(time.Now()) {
			creds = append(creds, Credentials{
				AccessKeyID:     cc.Credentials.AccessKeyId,
				SecretAccessKey: cc.Credentials.SecretAccessKey,
				SessionToken:    cc.Credentials.SessionToken,
			})
		}
	}
	return creds, nil
}

func StringifyCredentials(creds Credentials) (string, error) {
	b, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
