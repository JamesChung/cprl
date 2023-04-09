package remove

import (
	"fmt"

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
	shortMessage = "Remove branches"

	longMessage = templates.LongDesc(`
	Removes branches`)

	example = templates.Examples(`
	Remove branch with the aws-profile assigned to the default cprl profile:
	$ cprl codecommit branch remove

	Remove branch with a specified aws-profile:
	$ cprl codecommit branch remove --aws-profile=dev`)
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
		Use:     "remove",
		Aliases: []string{"r", "rem"},
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

	// Select a repository
	repo, err := cmd.Flags().GetString("repository")
	util.ExitOnErr(err)
	if repo == "" {
		repo, err = pterm.DefaultInteractiveSelect.
			WithDefaultText("Select a repository").
			WithOptions(repos).Show()
		util.ExitOnErr(err)
	}

	// Get branches for a given repository
	var branches []string
	util.Spinner("Retrieving branches...", func() {
		branches, err = ccClient.GetBranches(repo)
	})
	util.ExitOnErr(err)

	// Prompt for branches to remove
	branchSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(branches).Show("Select branches to remove")
	util.ExitOnErr(err)

	// Remove branches
	var results []util.Result[string]
	util.Spinner("Removing...", func() {
		results = util.DeleteBranches(ccClient, repo, branchSelections)
	})

	failCount := 0
	for _, r := range results {
		if r.Err != nil {
			pterm.Error.Printf("Failed to remove branch [%s]\n", r.Result)
			failCount++
			continue
		}
		pterm.Success.Printf("Successfully removed branch [%s]\n", r.Result)
	}

	if failCount > 0 {
		util.ExitOnErr(fmt.Errorf("Failed to remove %d branches\n", failCount))
	}
}
