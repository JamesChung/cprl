package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/console/open"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "AWS console"

	example = templates.Examples(`
	cprl console
	cprl console --aws-profile=dev-account`)
)

func consoleCommands() []*cobra.Command {
	return []*cobra.Command{
		open.NewCmd(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.Bool(
		"gov-cloud",
		false,
		"set context as gov-cloud",
	)
}

func NewCmdConsole() *cobra.Command {
	consoleCmd := &cobra.Command{
		Use:     "console",
		Aliases: []string{"co", "con"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(consoleCmd.PersistentFlags())
	util.AddGroup(consoleCmd, "Console:", consoleCommands()...)
	return consoleCmd
}
