package pr

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/codecommit/pr/approve"
	"github.com/JamesChung/cprl/cmd/codecommit/pr/closes"
	"github.com/JamesChung/cprl/cmd/codecommit/pr/create"
	"github.com/JamesChung/cprl/cmd/codecommit/pr/diff"
	"github.com/JamesChung/cprl/cmd/codecommit/pr/list"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Manage PRs"

	example = templates.Examples(`$ cprl codecommit pr`)
)

func prCommands() []*cobra.Command {
	return []*cobra.Command{
		approve.NewCmd(),
		closes.NewCmd(),
		create.NewCmd(),
		diff.NewCmd(),
		list.NewCmd(),
	}
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pr",
		Short:   shortMessage,
		Example: example,
	}
	util.AddGroup(cmd, "Commands:", prCommands()...)
	return cmd
}
