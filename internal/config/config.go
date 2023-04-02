package config

import (
	"fmt"
	"os"
	"path"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/JamesChung/cprl/pkg/util"
)

var (
	cprlConfigDir  = ".config/cprl"
	cprlConfigFile = fmt.Sprintf("%s/cprl.yaml", cprlConfigDir)
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
	cfgUserHomePath := path.Join(homeDir, cprlConfigDir)
	viper.AddConfigPath(cfgUserHomePath)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

type ConfigFile struct {
	Default ProfileBody `yaml:"default"`
}

type ProfileBody struct {
	Config   ProfileConfig   `yaml:"config"`
	Services ProfileServices `yaml:"services"`
}

type ProfileConfig struct {
	AWSProfile string `yaml:"aws-profile"`
}

type ProfileServices struct {
	CodeCommit ProfileServicesCodeCommit `yaml:"codecommit"`
}

type ProfileServicesCodeCommit struct {
	Repositories []string `yaml:"repositories"`
}

func Create() error {
	msg := "Would you like to create a cprl.yaml config file?"
	yes, err := pterm.DefaultInteractiveConfirm.Show(msg)
	if err != nil {
		return err
	}
	if !yes {
		return nil
	}
	cfg := ConfigFile{
		ProfileBody{
			Config: ProfileConfig{
				AWSProfile: "default",
			},
			Services: ProfileServices{
				CodeCommit: ProfileServicesCodeCommit{
					Repositories: []string{},
				},
			},
		},
	}
	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// Check if cprl config directory exists
	_, err = os.ReadDir(path.Join(homeDir, cprlConfigDir))
	if err != nil {
		if os.IsNotExist(err) {
			// If cprl config directory doesn't exist then create it
			err := os.Mkdir(path.Join(homeDir, cprlConfigDir), 0755)
			if err != nil {
				return err
			}
		}
	}
	// Create new cprl config file
	f, err := os.Create(path.Join(homeDir, cprlConfigFile))
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(yml)
	pterm.Success.Println("Created at", f.Name())
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
