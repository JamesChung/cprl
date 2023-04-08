package util

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"gopkg.in/ini.v1"

	"github.com/JamesChung/cprl/pkg/client"
)

type credentials struct {
	AccessKeyID     string `json:"sessionId"`
	SecretAccessKey string `json:"sessionKey"`
	SessionToken    string `json:"sessionToken"`
}

func GetCredentials(profile string) (aws.Credentials, error) {
	cfg, err := client.GetProfileConfig(profile)
	if err != nil {
		return aws.Credentials{}, err
	}
	ctx := context.Background()
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}
	return creds, nil
}

func WriteCredentials(section string, creds aws.Credentials) error {
	// Find user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	credentialsFile := path.Join(homeDir, ".aws/credentials")

	ini.PrettyFormat = false
	f, err := ini.Load(credentialsFile)
	if err != nil {
		return err
	}
	s := f.Section(section)
	s.NewKey("aws_access_key_id", creds.AccessKeyID)
	s.NewKey("aws_secret_access_key", creds.SecretAccessKey)
	s.NewKey("aws_session_token", creds.SessionToken)
	err = f.SaveTo(credentialsFile)
	if err != nil {
		return err
	}
	return nil
}

func StringifyCredentials(creds aws.Credentials) (string, error) {
	b, err := json.Marshal(credentials{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
	})
	if err != nil {
		return "", err
	}
	return string(b), nil
}
