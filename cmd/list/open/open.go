package open

import (
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	inutil "github.com/JamesChung/cprl/internal/util"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "List Open PRs"

	example = templates.Examples(`
	cprl list open
	cprl list open --profile=dev
	cprl list open --profile=dev --aws-profile=dev-account`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"author-arn",
		"",
		"filter by author",
	)
}

func NewCmdListOpen() *cobra.Command {
	openCmd := &cobra.Command{
		Use:     "open",
		Aliases: []string{"o"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(openCmd.Flags())
	return openCmd
}

func runCmd(cmd *cobra.Command, args []string) {
	inutil.Initialize(cmd)
	prs := inutil.GetPullRequests(cmd, types.PullRequestStatusEnumOpen)
	util.PRsToTable(prs)
}
