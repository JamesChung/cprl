package list

import (
	"github.com/pterm/pterm"
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

	longMessage = templates.LongDesc(`
	List a resource`)

	example = templates.Examples(`
	cprl codecommit pr list
	`)
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
	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(list.PersistentFlags())
	return list
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

	pterm.DefaultSpinner.Start("Getting PRs...")
	defer pterm.DefaultSpinner.Stop()
	prs := service.GetPullRequests(
		service.PullRequestInput{
			AuthorARN:    authorARN,
			Client:       ccClient,
			Repositories: repos,
			Status:       status,
		})
	service.PRsToTable(prs)
}
