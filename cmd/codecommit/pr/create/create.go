package create

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	shortMessage = "Create PRs"

	longMessage = templates.LongDesc(`
	Create a PR`)

	example = templates.Examples(`
	Create PR with the aws-profile assigned to the default cprl profile:
	$ cprl codecommit pr create

	Create PR with a specified aws-profile:
	$ cprl codecommit pr create --aws-profile=dev`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"repository",
		"",
		"repository name override",
	)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(cmd.Flags())
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

	// Get branches
	branches, err := util.Spinner("Getting branches...", func() ([]string, error) {
		branches, err := ccClient.GetBranches(cfg.Repository)
		if err != nil {
			return nil, err
		}
		return branches, nil
	})
	util.ExitOnErr(err)

	// Select source branch
	srcBranch, err := pterm.DefaultInteractiveSelect.
		WithDefaultText("Select a source branch").
		WithOptions(branches).Show()
	util.ExitOnErr(err)

	// Select destination branch
	destBranch, err := pterm.DefaultInteractiveSelect.
		WithDefaultText("Select a destination branch").
		WithOptions(branches).Show()
	util.ExitOnErr(err)

	// Input Title
	title, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("Input a title").Show()
	util.ExitOnErr(err)

	// Ask for description
	yes, err := pterm.DefaultInteractiveConfirm.
		WithDefaultText("Would you like a description?").Show()
	util.ExitOnErr(err)

	// Input Description
	var desc string
	if yes {
		desc, err = pterm.DefaultInteractiveTextInput.
			WithDefaultText("Input a Description").Show()
		util.ExitOnErr(err)
	}

	// Create PR
	msg, err := util.Spinner("Creating PR...", func() (string, error) {
		targets := []types.Target{
			{
				RepositoryName:       aws.String(cfg.Repository),
				SourceReference:      aws.String(srcBranch),
				DestinationReference: aws.String(destBranch),
			},
		}
		res, err := ccClient.CreatePR(targets, title, desc)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(
			"Created PR -> %s\n",
			aws.ToString(res.PullRequest.PullRequestId),
		), nil
	})
	util.ExitOnErr(err)
	pterm.Success.Println(msg)
}
