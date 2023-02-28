package closed

import (
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	inutil "github.com/JamesChung/cprl/internal/util"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "List Closed PRs"

	example = templates.Examples(`
	cprl list closed
	cprl list closed --profile=dev
	cprl list closed --profile=dev --aws-profile=dev-account`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"author-arn",
		"",
		"filter by author",
	)
}

func NewCmdListClosed() *cobra.Command {
	closedCmd := &cobra.Command{
		Use:     "closed",
		Aliases: []string{"c"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(closedCmd.PersistentFlags())
	return closedCmd
}

func runCmd(cmd *cobra.Command, args []string) {
	inutil.Initialize(cmd)
	prs := inutil.GetPullRequests(cmd, types.PullRequestStatusEnumClosed)
	util.PRsToTable(prs)
}
