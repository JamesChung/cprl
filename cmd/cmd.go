package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/JamesChung/cprl/cmd/codecommit"
	"github.com/JamesChung/cprl/pkg/util"
)

func cprlCommands() []*cobra.Command {
	return []*cobra.Command{
		codecommit.NewCmdCodeCommit(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"profile",
		"default",
		"profile name in cprl.yaml",
	)
	flags.String(
		"aws-profile",
		"",
		"AWS profile",
	)
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "cprl",
		Short:   "cprl",
		Version: "v0.1.0",
	}
	setPersistentFlags(rootCmd.PersistentFlags())
	util.AddGroup(rootCmd, "Services:", cprlCommands()...)
	return rootCmd
}
