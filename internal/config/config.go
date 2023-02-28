package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"

	"github.com/JamesChung/cprl/internal/cprlerrs"
)

type Config struct {
	Repositories []string
}

// Read config file cprl.yaml
func Read() error {
	// Set config file as cprl.yaml
	viper.SetConfigName("cprl")
	viper.SetConfigType("yaml")

	// Check working directory first
	viper.AddConfigPath(".")

	// Find user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Check user home config directory second
	cfgUserHomePath := path.Join(homeDir, ".config/cprl")
	viper.AddConfigPath(cfgUserHomePath)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func ProfileLookup(profile string) error {
	if !viper.InConfig(profile) {
		return fmt.Errorf("The profile [%s] does not exist", profile)
	}
	return nil
}

func GetRepositoriesConfig(profile string) ([]string, error) {
	repositories := viper.GetStringSlice(
		fmt.Sprintf("%s.repositories", profile),
	)
	if repositories == nil {
		return nil, cprlerrs.MissingRepositoriesErr
	}
	return repositories, nil
}

func GetConfig(profile string) (*Config, error) {
	if err := ProfileLookup("default"); err != nil {
		return nil, cprlerrs.MissingDefaultProfileErr
	}

	if err := ProfileLookup(profile); err != nil {
		return nil, err
	}

	repositories, err := GetRepositoriesConfig(profile)
	if err != nil {
		return nil, fmt.Errorf("[%s] profile in %w", profile, err)
	}

	return &Config{
		Repositories: repositories,
	}, nil
}
