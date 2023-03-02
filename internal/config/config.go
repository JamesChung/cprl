package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JamesChung/cprl/pkg/util"
)

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

func GetProfile(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "profile")
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetAWSProfile(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "aws-profile")
	if err != nil {
		return "", err
	}
	// If someone set the flag, use it because it takes precedent
	if str != "" {
		return str, nil
	}
	// Else we will look for the aws-profile in cprl.yaml
	profile, err := GetProfile(cmd)
	if err != nil {
		return "", err
	}
	if viper.IsSet(fmt.Sprintf("%s.config.aws-profile", profile)) {
		return viper.GetString(
			fmt.Sprintf(
				"%s.config.aws-profile",
				profile,
			),
		), nil
	}
	return "default", nil
}
