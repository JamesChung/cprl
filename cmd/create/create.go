package create

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/create/pr"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "create command"

	example = templates.Examples(`
	cprl create pr
	`)
)

func createCommands() []*cobra.Command {
	return []*cobra.Command{
		pr.NewCmdCreatePR(),
	}
}

func NewCmdCreate() *cobra.Command {
	create := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   shortMessage,
		Example: example,
	}
	util.AddGroup(create, "Commands:", createCommands()...)
	return create
}
