package codecommit

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/codecommit/branch"
	"github.com/JamesChung/cprl/cmd/codecommit/pr"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "CodeCommit"

	example = templates.Examples(`$ cprl codecommit`)
)

func codeCommitCommands() []*cobra.Command {
	return []*cobra.Command{
		branch.NewCmd(),
		pr.NewCmd(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	// TODO Possibly revisit and add flags if needed
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "codecommit",
		Aliases: []string{"cc"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(cmd.PersistentFlags())
	util.AddGroup(cmd, "Commands:", codeCommitCommands()...)
	return cmd
}
