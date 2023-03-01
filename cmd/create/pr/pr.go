package pr

import (
	"github.com/aws/aws-sdk-go-v2/aws"
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
	shortMessage = "Create PRs"

	example = templates.Examples(`
	cprl create pr
	cprl create pr --aws-profile=dev-account`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"repository",
		"",
		"repository name",
	)
}

func NewCmdCreatePR() *cobra.Command {
	createPRCmd := &cobra.Command{
		Use:     "pr",
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(createPRCmd.Flags())
	return createPRCmd
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

	// Select a repository
	repo, err := pterm.DefaultInteractiveSelect.WithDefaultText(
		"Select a repository",
	).WithOptions(
		repos,
	).Show()
	util.ExitOnErr(err)

	// Get branches
	branches, err := ccClient.GetBranches(repo)
	util.ExitOnErr(err)

	// Select source srcBranch
	srcBranch, err := pterm.DefaultInteractiveSelect.WithDefaultText(
		"Select a source branch",
	).WithOptions(
		branches,
	).Show()
	util.ExitOnErr(err)

	// Select destination branch
	destBranch, err := pterm.DefaultInteractiveSelect.WithDefaultText(
		"Select a destination branch",
	).WithOptions(
		branches,
	).Show()
	util.ExitOnErr(err)

	// Input Title
	title, err := pterm.DefaultInteractiveTextInput.WithDefaultText(
		"Input a title").Show()
	util.ExitOnErr(err)

	// Ask for description
	yes, err := pterm.DefaultInteractiveConfirm.WithDefaultText(
		"Would you like a description?").Show()
	util.ExitOnErr(err)

	// Input Description
	var desc string
	if yes {
		desc, err = pterm.DefaultInteractiveTextInput.WithDefaultText(
			"Input a Description").Show()
		util.ExitOnErr(err)
	}

	// Create PR
	targets := []types.Target{
		{
			RepositoryName:       aws.String(repo),
			SourceReference:      aws.String(srcBranch),
			DestinationReference: aws.String(destBranch),
		},
	}
	_, err = ccClient.CreatePR(targets, title, desc)
	util.ExitOnErr(err)
}
