package diff

import (
	"bytes"
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
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Diff PRs"

	longMessage = templates.LongDesc(`
	Diff a PR`)

	example = templates.Examples(`
	Diff PR with the aws-profile assigned to the default cprl profile:
	$ cprl codecommit pr diff

	Diff PR with a specified aws-profile:
	$ cprl codecommit pr diff --aws-profile=dev`)
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
		Use:     "diff",
		Aliases: []string{"d"},
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
	var prs []*codecommit.GetPullRequestOutput
	util.Spinner("Retrieving PRs...", func() {
		prs, err = ccClient.GetPullRequests(client.PullRequestInput{
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

	// Get diff metadata info between targets from CodeCommit
	var diffOut []*codecommit.GetDifferencesOutput
	util.Spinner("Getting Differences...", func() {
		for _, t := range prMap[prSelection].PullRequest.PullRequestTargets {
			diffOut, err = ccClient.GetDifferences(
				aws.String(repo),
				t.DestinationReference,
				t.SourceReference,
			)
		}
	})
	util.ExitOnErr(err)

	// Concurrently generate individual file diffs
	var results, badResults []client.Result[[]byte]
	util.Spinner("Generating Differences...", func() {
		results = ccClient.GenerateDiffs(repo, diffOut)
	})

	// Prompt user for the name of the diff file
	dFileName, err := pterm.DefaultInteractiveTextInput.Show("Submit name of diff file")
	util.ExitOnErr(err)

	// Create diff file
	f, err := os.Create(dFileName)
	util.ExitOnErr(err)
	defer f.Close()

	// Use a bytes buffer to write the complete diff file from individual file diff bytes
	buf := bytes.Buffer{}
	for _, res := range results {
		if res.Err != nil {
			badResults = append(badResults, res)
			continue
		}
		_, err = buf.Write(res.Result)
		if err != nil {
			util.ExitOnErr(err)
		}
	}

	// Write diff to file
	_, err = f.Write(buf.Bytes())
	if err != nil {
		util.ExitOnErr(err)
	}

	// End if there are no errors
	if len(badResults) == 0 {
		return
	}

	// Notify user errors have occurred, prompt if they would like a report
	errReport, err := pterm.DefaultInteractiveConfirm.Show("There were some errors would you like a report?")
	if err != nil {
		util.ExitOnErr(err)
	}

	// End if the user would not like an error report
	if !errReport {
		return
	}

	// Prompt user for error log file name
	errFileName, err := pterm.DefaultInteractiveTextInput.Show("Submit name of error log file")
	if err != nil {
		util.ExitOnErr(err)
	}

	// Create error log file
	e, err := os.Create(errFileName)
	util.ExitOnErr(err)
	defer e.Close()

	errBuf := bytes.Buffer{}
	for _, res := range badResults {
		res.Result = append(res.Result, []byte(fmt.Sprintf(":\n%s\n\n", res.Err.Error()))...)
		errBuf.Write(res.Result)
	}

	// Write errors to file
	_, err = e.Write(errBuf.Bytes())
	if err != nil {
		util.ExitOnErr(err)
	}
}
