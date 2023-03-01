package pr

import (
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
	shortMessage = "List PRs"

	example = templates.Examples(`
	cprl list pr
	cprl list pr --profile=dev
	cprl list pr --profile=dev --aws-profile=dev-account`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"author-arn",
		"",
		"filter by author",
	)
	flags.Bool(
		"closed",
		false,
		"filter closed PRs",
	)
}

func NewCmdListPR() *cobra.Command {
	openCmd := &cobra.Command{
		Use:     "pr",
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(openCmd.Flags())
	return openCmd
}

func runCmd(cmd *cobra.Command, args []string) {
	profile, err := config.GetProfile(cmd)
	util.ExitOnErr(err)

	awsProfile, err := config.GetAWSProfile(cmd)
	util.ExitOnErr(err)

	ccClient, err := client.NewCodeCommitClient(awsProfile)
	util.ExitOnErr(err)

	repos, err := cc.GetRepositories(profile)
	util.ExitOnErr(err)

	authorARN, err := cc.GetAuthorARN(cmd)
	util.ExitOnErr(err)

	status, err := cc.GetClosed(cmd)
	util.ExitOnErr(err)

	prs := service.GetPullRequests(
		service.PullRequestInput{
			AuthorARN:    authorARN,
			Client:       ccClient,
			Repositories: repos,
			Status:       status,
		})
	service.PRsToTable(prs)
}
