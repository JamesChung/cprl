package branch

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/codecommit/branch/del"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Manage branches"

	example = templates.Examples(`
	cprl codecommit branch
	cprl codecommit branch --aws-profile=dev-account`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	// TODO
}

func branchCommands() []*cobra.Command {
	return []*cobra.Command{
		del.NewCmd(),
	}
}

func NewCmd() *cobra.Command {
	prCmd := &cobra.Command{
		Use:     "branch",
		Aliases: []string{"br"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(prCmd.Flags())
	util.AddGroup(prCmd, "Commands:", branchCommands()...)
	return prCmd
}
