package approve

import (
	"context"
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
	cprl codecommit pr approve
	`)
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

func NewCmdApprovePR() *cobra.Command {
	approve := &cobra.Command{
		Use:     "approve",
		Aliases: []string{"a"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(approve.PersistentFlags())
	return approve
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
	var prs []*codecommit.GetPullRequestOutput
	util.Spinner("Retrieving PRs...", func() {
		prs, err = util.GetPullRequestsSlice(util.PullRequestInput{
			Client:       ccClient,
			AuthorARN:    authorARN,
			Repositories: []string{repo},
			Status:       types.PullRequestStatusEnumOpen,
		})
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

	// Prompt for PRs to approve
	prSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(li).Show("Select PRs to approve")
	util.ExitOnErr(err)

	// Approve PRs
	util.SpinnerWithStatusMsg("Approving...", func() (string, error) {
		ctx := context.Background()
		for _, v := range prSelections {
			_, err = ccClient.Client.UpdatePullRequestApprovalState(
				ctx, &codecommit.UpdatePullRequestApprovalStateInput{
					ApprovalState: types.ApprovalStateApprove,
					PullRequestId: prMap[v].PullRequest.PullRequestId,
					RevisionId:    prMap[v].PullRequest.RevisionId,
				})
			if err != nil {
				// Return and let outer scope handle error
				return "", nil
			}
		}
		return "Successfully approved!", nil
	})
	util.ExitOnErr(err)
}
