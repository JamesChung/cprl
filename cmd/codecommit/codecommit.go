package codecommit

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/cmd/codecommit/pr"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "CodeCommit"

	example = templates.Examples(`
	cprl codecommit pr
	cprl codecommit pr --aws-profile=dev-account`)
)

func codeCommitCommands() []*cobra.Command {
	return []*cobra.Command{
		pr.NewCmdPR(),
	}
}

func setPersistentFlags(flags *pflag.FlagSet) {
	// flags.String(
	// 	"author-arn",
	// 	"",
	// 	"",
	// )
}

func NewCmdCodeCommit() *cobra.Command {
	codeCommitCmd := &cobra.Command{
		Use:     "codecommit",
		Aliases: []string{"cc"},
		Short:   shortMessage,
		Example: example,
	}
	setPersistentFlags(codeCommitCmd.PersistentFlags())
	util.AddGroup(codeCommitCmd, "Commands:", codeCommitCommands()...)
	return codeCommitCmd
}
