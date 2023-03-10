package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/JamesChung/cprl/cmd/codecommit"
	"github.com/JamesChung/cprl/cmd/console"
	"github.com/JamesChung/cprl/pkg/util"
)

func cprlCommands() []*cobra.Command {
	return []*cobra.Command{
		codecommit.NewCmdCodeCommit(),
		console.NewCmdConsole(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"profile",
		"default",
		"references a profile in cprl.yaml",
	)
	flags.String(
		"aws-profile",
		"",
		"overrides [aws-profile] value in cprl.yaml",
	)
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "cprl",
		Short:   "cprl",
		Version: "v0.0.0-alpha",
	}
	setPersistentFlags(rootCmd.PersistentFlags())
	util.AddGroup(rootCmd, "Services:", cprlCommands()...)
	return rootCmd
}
