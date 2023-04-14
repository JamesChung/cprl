package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConsoleConfig struct {
	*Config
	Interactive bool
	IsGovCloud  bool
}

func NewConsoleConfig(cmd *cobra.Command) (*ConsoleConfig, error) {
	c, err := NewConfig(cmd)
	if err != nil {
		return nil, err
	}
	cfg := &ConsoleConfig{}
	cfg.Config = c
	interactive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		return nil, err
	}
	cfg.Interactive = interactive
	isGovCloud, err := IsGovCloud(cmd, cfg.Profile)
	if err != nil {
		return nil, err
	}
	cfg.IsGovCloud = isGovCloud
	return cfg, nil
}

func IsGovCloud(cmd *cobra.Command, profile string) (bool, error) {
	gov, err := cmd.Flags().GetBool("gov-cloud")
	if err != nil {
		return gov, err
	}
	if gov {
		return gov, nil
	}
	gov = viper.GetBool(
		fmt.Sprintf(
			"%s.services.console.gov-cloud",
			profile,
		),
	)
	return gov, nil
}
