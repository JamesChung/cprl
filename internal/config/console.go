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
	cfg := &ConsoleConfig{Config: c}
	cfg.Interactive, _ = cmd.Flags().GetBool("interactive")
	cfg.IsGovCloud = IsGovCloud(cmd, cfg.Profile)
	return cfg, nil
}

func IsGovCloud(cmd *cobra.Command, profile string) bool {
	gov, _ := cmd.Flags().GetBool("gov-cloud")
	if gov {
		return gov
	}

	// Fallback to cprl.yaml profile value if flag is not set
	gov = viper.GetBool(
		fmt.Sprintf(
			"%s.services.console.gov-cloud",
			profile,
		),
	)
	return gov
}
