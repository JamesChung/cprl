package pr

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/codecommit/pr/create"
	"github.com/JamesChung/cprl/cmd/codecommit/pr/list"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Manage PRs"

	example = templates.Examples(`
	cprl codecommit pr
	cprl codecommit pr --aws-profile=dev-account`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	// flags.String(
	// 	"repository",
	// 	"",
	// 	"repository name",
	// )
}

func prCommands() []*cobra.Command {
	return []*cobra.Command{
		create.NewCmdCreatePR(),
		list.NewCmdListPR(),
	}
}

func NewCmdPR() *cobra.Command {
	prCmd := &cobra.Command{
		Use:     "pr",
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(prCmd.Flags())
	util.AddGroup(prCmd, "Commands:", prCommands()...)
	return prCmd
}
