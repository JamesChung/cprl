package del

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
	shortMessage = "Delete branches"

	longMessage = templates.LongDesc(`
	Delete a branch`)

	example = templates.Examples(`
	cprl codecommit branch delete
	`)
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
		Use:     "delete",
		Aliases: []string{"d", "del"},
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

	// Prompt for branches to delete
	branchSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(branches).Show("Select branches to delete")
	util.ExitOnErr(err)

	// Delete branches
	var results []util.Result[string]
	util.Spinner("Deleting...", func() {
		results = util.DeleteBranches(ccClient, repo, branchSelections)
	})

	failCount := 0
	for _, r := range results {
		if r.Err != nil {
			pterm.Error.Printf("Failed to delete branch [%s]\n", r.Result)
			failCount++
			continue
		}
		pterm.Success.Printf("Successfully deleted branch [%s]\n", r.Result)
	}

	if failCount > 0 {
		util.ExitOnErr(fmt.Errorf("Failed to delete %d branches\n", failCount))
	}
}
