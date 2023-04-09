package cloudwatch

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	// "github.com/JamesChung/cprl/cmd/cloudwatch/logs"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "CloudWatch"

	example = templates.Examples(`$ cprl cloudwatch`)
)

func cloudwatchCommands() []*cobra.Command {
	return []*cobra.Command{
		// logs.NewCmd(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	// TODO:
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloudwatch",
		Aliases: []string{"cw"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(cmd.PersistentFlags())
	util.AddGroup(cmd, "Commands:", cloudwatchCommands()...)
	return cmd
}
