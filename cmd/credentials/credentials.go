package credentials

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/credentials/assume"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "AWS credentials"

	example = templates.Examples(`
	cprl credentials
	cprl credentials --aws-profile=dev-account`)
)

func credentialsCommands() []*cobra.Command {
	return []*cobra.Command{
		assume.NewCmd(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.Bool(
		"gov-cloud",
		false,
		"set context as gov-cloud",
	)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credentials",
		Aliases: []string{"cr", "creds"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(cmd.PersistentFlags())
	util.AddGroup(cmd, "Credentials:", credentialsCommands()...)
	return cmd
}
