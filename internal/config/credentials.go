package config

import "github.com/spf13/cobra"

type CredentialsConfig struct {
	*Config
	RoleARN       string
	SessionName   string
	OutputProfile string
	OutputStyle   string
	IsJSON        bool
	Interactive   bool
}

func NewCredentialsConfig(cmd *cobra.Command) (*CredentialsConfig, error) {
	c, err := NewConfig(cmd)
	if err != nil {
		return nil, err
	}
	cfg := &CredentialsConfig{Config: c}
	cfg.RoleARN, _ = cmd.Flags().GetString("role-arn")
	cfg.SessionName, _ = cmd.Flags().GetString("session-name")
	cfg.OutputProfile, _ = cmd.Flags().GetString("output-profile")
	cfg.OutputStyle, _ = cmd.Flags().GetString("output-style")
	cfg.IsJSON, _ = cmd.Flags().GetBool("json")
	cfg.Interactive, _ = cmd.Flags().GetBool("interactive")
	return cfg, nil
}
