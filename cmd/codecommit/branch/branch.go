package branch

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/codecommit/branch/remove"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Manage branches"

	example = templates.Examples(`$ cprl codecommit branch`)
)

func branchCommands() []*cobra.Command {
	return []*cobra.Command{
		remove.NewCmd(),
	}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "branch",
		Aliases: []string{"br"},
		Short:   shortMessage,
		Example: example,
	}
	util.AddGroup(cmd, "Commands:", branchCommands()...)
	return cmd
}
