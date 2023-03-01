package closed

import (
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	cc "github.com/JamesChung/cprl/internal/config/services/codecommit"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/service"
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
	profile, err := config.GetProfileFlag(cmd)
	util.ExitOnErr(err)

	ccClient, err := client.NewCodeCommitClient(profile)
	util.ExitOnErr(err)

	repos, err := cc.GetRepositories(profile)
	util.ExitOnErr(err)

	authorARN, err := cc.GetAuthorARNFlag(cmd)
	util.ExitOnErr(err)

	prs := service.GetPullRequests(
		service.PullRequestInput{
			AuthorARN:    authorARN,
			Client:       ccClient,
			Repositories: repos,
			Status:       types.PullRequestStatusEnumClosed,
		})
	service.PRsToTable(prs)
}
