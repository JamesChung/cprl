package list

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/list/closed"
	"github.com/JamesChung/cprl/cmd/list/open"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "List PRs"

	example = templates.Examples(`
	cprl list open
	cprl list closed
	`)
)

func listCommands() []*cobra.Command {
	return []*cobra.Command{
		open.NewCmdListOpen(),
		closed.NewCmdListClosed(),
	}
}

func NewCmdList() *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   shortMessage,
		Example: example,
	}
	util.AddGroup(list, "Commands:", listCommands()...)
	return list
}
