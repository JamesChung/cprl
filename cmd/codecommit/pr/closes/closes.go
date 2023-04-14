package closes

import (
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
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Close PRs"

	longMessage = templates.LongDesc(`
	Close a PR`)

	example = templates.Examples(`
	Close PR with the aws-profile assigned to the default cprl profile:
	$ cprl codecommit pr close

	Close PR with a specified aws-profile:
	$ cprl codecommit pr close --aws-profile=dev`)
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
		Use:     "close",
		Aliases: []string{"cl"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(cmd.PersistentFlags())
	return cmd
}

func runCmd(cmd *cobra.Command, args []string) {
	cfg, err := config.NewCodeCommitConfig(cmd)
	util.ExitOnErr(err)
	repos, err := config.GetRepositories(cfg.Profile)
	util.ExitOnErr(err)
	ccClient, err := client.NewCodeCommitClient(cfg.AWSProfile)
	util.ExitOnErr(err)

	// Select a repository
	if cfg.Repository == "" {
		cfg.Repository, err = pterm.DefaultInteractiveSelect.
			WithOptions(repos).Show("Select a repository")
		util.ExitOnErr(err)
	}

	// Get PRs for a given repository
	prs, err := util.Spinner("Retrieving PRs...", func() ([]*codecommit.GetPullRequestOutput, error) {
		prs, err := ccClient.GetPullRequests(client.PullRequestInput{
			AuthorARN:    cfg.AuthorARN,
			Repositories: []string{cfg.Repository},
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
		util.ExitOnErr(fmt.Errorf("[%s] has no available PRs", cfg.Repository))
	}

	// Prompt for PRs to close
	prSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(li).Show("Select PRs to close")
	util.ExitOnErr(err)

	// Close PRs
	res, _ := util.Spinner("Closing...", func() ([]client.Result[string], error) {
		res := ccClient.ClosePRs(prMap, prSelections)
		return res, nil
	})

	// Output PR close results
	errCount := 0
	for _, r := range res {
		if r.Err != nil {
			pterm.Error.Printf("Failed to close PR [%s]: %s\n", r.Result, r.Err)
			errCount++
			continue
		}
		pterm.Success.Printf("Successfully closed PR [%s]\n", r.Result)
	}
	if errCount > 0 {
		os.Exit(1)
	}
}
