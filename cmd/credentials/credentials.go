package credentials

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/credentials/assume"
	"github.com/JamesChung/cprl/cmd/credentials/list"
	"github.com/JamesChung/cprl/cmd/credentials/output"
	"github.com/JamesChung/cprl/cmd/credentials/remove"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "AWS credentials"

	example = templates.Examples(`$ cprl credentials`)
)

func credentialsCommands() []*cobra.Command {
	return []*cobra.Command{
		assume.NewCmd(),
		list.NewCmd(),
		output.NewCmd(),
		remove.NewCmd(),
	}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credentials",
		Aliases: []string{"cr", "creds"},
		Short:   shortMessage,
		Example: example,
	}
	util.AddGroup(cmd, "Commands:", credentialsCommands()...)
	return cmd
}
