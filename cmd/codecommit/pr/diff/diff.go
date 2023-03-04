package diff

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/pterm/pterm"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	cc "github.com/JamesChung/cprl/internal/config/services/codecommit"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Diff PRs"

	longMessage = templates.LongDesc(`
	Diff a PR`)

	example = templates.Examples(`
	cprl codecommit pr diff
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

func NewCmdDiffPR() *cobra.Command {
	diff := &cobra.Command{
		Use:     "diff",
		Aliases: []string{"d"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(diff.PersistentFlags())
	return diff
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
		prs, err = util.GetPullRequests(util.PullRequestInput{
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
	prSelection, err := pterm.DefaultInteractiveSelect.
		WithOptions(li).Show("Select PR to diff")
	util.ExitOnErr(err)

	var diffOut []*codecommit.GetDifferencesOutput
	util.Spinner("Getting Differences...", func() {
		for _, t := range prMap[prSelection].PullRequest.PullRequestTargets {
			diffOut, err = ccClient.GetDifferences(
				aws.String(repo),
				t.SourceReference,
				t.DestinationReference,
			)
		}
	})
	util.ExitOnErr(err)

	for _, do := range diffOut {
		for _, d := range do.Differences {
			pterm.Info.Println(d.ChangeType, aws.ToString(d.BeforeBlob.Path))
			pterm.Println("-----------------------------------")
			bob, err := ccClient.Client.GetBlob(context.TODO(), &codecommit.GetBlobInput{
				BlobId:         d.BeforeBlob.BlobId,
				RepositoryName: aws.String(repo),
			})
			util.ExitOnErr(err)
			boa, err := ccClient.Client.GetBlob(context.TODO(), &codecommit.GetBlobInput{
				BlobId:         d.AfterBlob.BlobId,
				RepositoryName: aws.String(repo),
			})
			util.ExitOnErr(err)
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(
				string(bob.Content),
				string(boa.Content),
				false,
			)
			cmd.Println(dmp.DiffPrettyText(diffs))
		}
	}
}
