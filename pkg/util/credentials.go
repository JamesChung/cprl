package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"gopkg.in/ini.v1"

	"github.com/JamesChung/cprl/pkg/client"
)

type sessionCredentials struct {
	SessionID    string `json:"sessionId"`
	SessionKey   string `json:"sessionKey"`
	SessionToken string `json:"sessionToken"`
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

func GetAllProfiles() ([]string, error) {
	configProfiles, err := GetConfigProfiles()
	if err != nil {
		return nil, err
	}
	credentialsProfiles, err := GetCredentialsProfiles()
	if err != nil {
		return nil, err
	}
	pMap := make(map[string]struct{}, len(configProfiles)+len(credentialsProfiles))
	profiles := make([]string, 0, len(configProfiles)+len(credentialsProfiles))
	for _, p := range configProfiles {
		pMap[p] = struct{}{}
	}
	for _, p := range credentialsProfiles {
		pMap[p] = struct{}{}
	}
	for k := range pMap {
		profiles = append(profiles, k)
	}
	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles were found in config and credentials")
	}
	return profiles, nil
}

func GetConfigProfiles() ([]string, error) {
	configFile, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	ini.PrettyFormat = false
	f, err := ini.Load(configFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	return f.SectionStrings()[1:], nil
}

func GetCredentialsProfiles() ([]string, error) {
	credentialsFile, err := getCredentialsFilePath()
	if err != nil {
		return nil, err
	}
	ini.PrettyFormat = false
	f, err := ini.Load(credentialsFile)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	return f.SectionStrings()[1:], nil
}

func getConfigFilePath() (string, error) {
	// Find user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configFile := path.Join(homeDir, ".aws/config")
	return configFile, nil
}

func getCredentialsFilePath() (string, error) {
	// Find user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	credentialsFile := path.Join(homeDir, ".aws/credentials")
	return credentialsFile, nil
}

func WriteCredentials(section string, creds aws.Credentials) error {
	credentialsFile, err := getCredentialsFilePath()
	if err != nil {
		return err
	}

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

func RemoveProfiles(profiles []string) error {
	credentialsFile, err := getCredentialsFilePath()
	if err != nil {
		return err
	}
	ini.PrettyFormat = false
	f, err := ini.Load(credentialsFile)
	if err != nil {
		return err
	}
	for _, profile := range profiles {
		f.DeleteSection(profile)
	}
	err = f.SaveTo(credentialsFile)
	if err != nil {
		return err
	}
	return nil
}

type credentials struct {
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
	Source          string `json:"Source"`
}

func StringifyCredentials(creds aws.Credentials) (string, error) {
	var b []byte
	var err error
	if creds.CanExpire {
		b, err = json.Marshal(creds)
		if err != nil {
			return "", err
		}
	} else {
		b, err = json.Marshal(credentials{
			AccessKeyID:     creds.AccessKeyID,
			SecretAccessKey: creds.SecretAccessKey,
			SessionToken:    creds.SessionToken,
			Source:          creds.Source,
		})
		if err != nil {
			return "", err
		}
	}
	return string(b), nil
}

func StringifySessionCredentials(creds aws.Credentials) (string, error) {
	b, err := json.Marshal(sessionCredentials{
		SessionID:    creds.AccessKeyID,
		SessionKey:   creds.SecretAccessKey,
		SessionToken: creds.SessionToken,
	})
	if err != nil {
		return "", err
	}
	return string(b), nil
}
