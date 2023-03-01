package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/JamesChung/cprl/cmd/create"
	"github.com/JamesChung/cprl/cmd/list"
	"github.com/JamesChung/cprl/pkg/util"
)

func cprlCommands() []*cobra.Command {
	return []*cobra.Command{
		create.NewCmdCreate(),
		list.NewCmdList(),
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
		"default",
		"AWS profile",
	)
}

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cprl",
		Short: "cprl",
	}
	setPersistentFlags(rootCmd.PersistentFlags())
	util.AddGroup(rootCmd, "Commands:", cprlCommands()...)
	return rootCmd
}
