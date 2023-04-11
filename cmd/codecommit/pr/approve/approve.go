package approve

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
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
	shortMessage = "Approve PRs"

	longMessage = templates.LongDesc(`
	Approve a PR`)

	example = templates.Examples(`
	Approve PR with the aws-profile assigned to the default cprl profile:
	$ cprl codecommit pr approve

	Approve PR with a specified aws-profile:
	$ cprl codecommit pr approve --aws-profile=dev`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"author-arn",
		"",
		"filter by author",
	)
	flags.String(
		"repository",
		"",
		"repository name override",
	)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "approve",
		Aliases: []string{"a"},
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
	repos, err := cc.GetRepositories(profile)
	util.ExitOnErr(err)
	awsProfile, err := config.GetAWSProfile(cmd)
	util.ExitOnErr(err)
	ccClient, err := client.NewCodeCommitClient(awsProfile)
	util.ExitOnErr(err)
	authorARN, err := cc.GetAuthorARN(cmd)
	util.ExitOnErr(err)

	// Select a repository
	repo, err := cmd.Flags().GetString("repository")
	util.ExitOnErr(err)
	if repo == "" {
		repo, err = pterm.DefaultInteractiveSelect.
			WithDefaultText("Select a repository").
			WithOptions(repos).Show()
		util.ExitOnErr(err)
	}

	// Get PRs for a given repository
	prs, err := util.Spinner("Retrieving PRs...", func() ([]*codecommit.GetPullRequestOutput, error) {
		prs, err := ccClient.GetPullRequests(client.PullRequestInput{
			AuthorARN:    authorARN,
			Repositories: []string{repo},
			Status:       types.PullRequestStatusEnumOpen,
		})
		if err != nil {
			return nil, err
		}
		return prs, nil
	})
	util.ExitOnErr(err)

	// Construct human readable selection options and map the associated option
	// to the codecommit responses
	prMap := make(map[string]*codecommit.GetPullRequestOutput)
	li := []string{}
	for i, p := range prs {
		s := fmt.Sprintf("%s: %s",
			aws.ToString(p.PullRequest.PullRequestId),
			aws.ToString(p.PullRequest.Title))
		li = append(li, s)
		prMap[s] = prs[i]
	}
	if len(li) == 0 {
		util.ExitOnErr(fmt.Errorf("[%s] has no available PRs", repo))
	}

	// Prompt for PRs to approve
	prSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(li).Show("Select PRs to approve")
	util.ExitOnErr(err)

	// Approve PRs
	res, _ := util.Spinner("Approving...", func() ([]client.Result[string], error) {
		res := ccClient.ApprovePRs(prMap, prSelections)
		return res, nil
	})

	// Output PR approval results
	errCount := 0
	for _, r := range res {
		if r.Err != nil {
			pterm.Error.Printf("Failed to approve PR [%s]: %s\n", r.Result, r.Err)
			errCount++
			continue
		}
		pterm.Success.Printf("Successfully approved PR [%s]\n", r.Result)
	}
	if errCount > 0 {
		util.ExitOnErr(fmt.Errorf("%d PRs failed to be approved\n", errCount))
	}
}
