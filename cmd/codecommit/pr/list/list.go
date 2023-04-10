package list

import (
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	cc "github.com/JamesChung/cprl/internal/config/services/codecommit"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "List PRs"

	longMessage = templates.LongDesc(`
	List a resource`)

	example = templates.Examples(`
	List PR with the aws-profile assigned to the default cprl profile:
	$ cprl codecommit pr list

	List PR with a specified aws-profile:
	$ cprl codecommit pr list --aws-profile=dev`)
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(cmd.PersistentFlags())
	return cmd
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

	// Get repository query list
	repoSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(repos).Show("Select repositories")
	util.ExitOnErr(err)

	// Get table headers
	tblSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions([]string{
			"Repository",
			"Author",
			"ID",
			"Title",
			"Source",
			"Destination",
			"CreationDate",
			"LastActivityDate",
		}).Show("Select table headers")
	util.ExitOnErr(err)

	// Get PR IDs
	var prIDs [][]string
	util.Spinner("Getting PR IDs...", func() {
		prIDs, err = util.GetPullRequestIDs(
			util.PullRequestInput{
				AuthorARN:    authorARN,
				Client:       ccClient,
				Repositories: repoSelections,
				Status:       status,
			})
	})
	util.ExitOnErr(err)

	// Get PR information
	var prInfoList []*codecommit.GetPullRequestOutput
	util.Spinner("Getting PR Information...", func() {
		prInfoList, err = util.GetPullRequestInfoFromIDs(ccClient, prIDs)
	})
	util.ExitOnErr(err)

	// Generate table
	var tbl *pterm.TablePrinter
	util.Spinner("Generating Table...", func() {
		tbl = util.PRsToTable(util.GenerateTableHeaders(tblSelections), prInfoList)
	})
	err = tbl.Render()
	util.ExitOnErr(err)
}
