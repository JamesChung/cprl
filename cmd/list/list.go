package list

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/list/pr"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "list command"

	longMessage = templates.LongDesc(`
	List a resource`)

	example = templates.Examples(`
	cprl list pr
	`)
)

func listCommands() []*cobra.Command {
	return []*cobra.Command{
		pr.NewCmdListPR(),
	}
}

func NewCmdList() *cobra.Command {
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
	}
	util.AddGroup(list, "Commands:", listCommands()...)
	return list
}
