package util

import (
	"context"
	"encoding/json"

	"github.com/JamesChung/cprl/pkg/client"
	"github.com/aws/aws-sdk-go-v2/aws"
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
