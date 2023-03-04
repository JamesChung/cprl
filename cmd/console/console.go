package console

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "AWS console"

	example = templates.Examples(`
	cprl console
	cprl console --aws-profile=dev-account`)
)

func consoleCommands() []*cobra.Command {
	return []*cobra.Command{}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	// TODO
}

func NewCmdConsole() *cobra.Command {
	consoleCmd := &cobra.Command{
		Use:     "console",
		Aliases: []string{"con"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(consoleCmd.PersistentFlags())
	util.AddGroup(consoleCmd, "Commands:", consoleCommands()...)
	return consoleCmd
}

func runCmd(cmd *cobra.Command, args []string) {
	cmd.Println("This command is under development")
}
