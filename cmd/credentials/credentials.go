package credentials

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/credentials/assume"
	"github.com/JamesChung/cprl/cmd/credentials/clear"
	"github.com/JamesChung/cprl/cmd/credentials/list"
	"github.com/JamesChung/cprl/cmd/credentials/output"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "AWS credentials"

	example = templates.Examples(`$ cprl credentials`)
)

func credentialsCommands() []*cobra.Command {
	return []*cobra.Command{
		assume.NewCmd(),
		clear.NewCmd(),
		list.NewCmd(),
		output.NewCmd(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	// TODO Possibly revisit and add flags if needed
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credentials",
		Aliases: []string{"cr", "creds"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(cmd.PersistentFlags())
	util.AddGroup(cmd, "Commands:", credentialsCommands()...)
	return cmd
}
