package logs

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "CloudWatch Logs"

	example = templates.Examples(`$ cprl cloudwatch logs`)
)

func cloudwatchLogsCommands() []*cobra.Command {
	return []*cobra.Command{}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs",
		Aliases: []string{"l", "log"},
		Short:   shortMessage,
		Example: example,
	}
	util.AddGroup(cmd, "Commands:", cloudwatchLogsCommands()...)
	return cmd
}
