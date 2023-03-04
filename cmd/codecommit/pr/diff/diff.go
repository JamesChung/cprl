package diff

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	cc "github.com/JamesChung/cprl/internal/config/services/codecommit"
	"github.com/JamesChung/cprl/internal/diff"
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

	var diffResult []byte
	util.Spinner("Generating Differences...", func() {
		for _, do := range diffOut {
			for _, d := range do.Differences {
				bob, err := ccClient.Client.GetBlob(context.TODO(), &codecommit.GetBlobInput{
					BlobId:         d.BeforeBlob.BlobId,
					RepositoryName: aws.String(repo),
				})
				// Let outer scope handle error
				if err != nil {
					return
				}
				boa, err := ccClient.Client.GetBlob(context.TODO(), &codecommit.GetBlobInput{
					BlobId:         d.AfterBlob.BlobId,
					RepositoryName: aws.String(repo),
				})
				if err != nil {
					// Let outer scope handle error
					return
				}
				diffResult = diff.Diff(
					aws.ToString(d.BeforeBlob.Path),
					bob.Content, aws.ToString(d.AfterBlob.Path),
					boa.Content,
				)
			}
		}
	})
	util.ExitOnErr(err)

	// Prompt user if they want a diff file
	ds, err := pterm.DefaultInteractiveConfirm.Show("Would you like to output the diff to a file?")
	util.ExitOnErr(err)
	// Print diff and exit early if user doesn't want a diff file
	if !ds {
		cmd.Println(string(diffResult))
		return
	}
	// Prompt user for the name of the diff file
	dFileName, err := pterm.DefaultInteractiveTextInput.Show("Submit name of diff file")
	util.ExitOnErr(err)
	// Create diff file
	f, err := os.Create(dFileName)
	util.ExitOnErr(err)
	defer f.Close()
	// Write diff to file
	f.Write(diffResult)
}
