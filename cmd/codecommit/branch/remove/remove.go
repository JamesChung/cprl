package remove

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
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
		Aliases: []string{"r", "rm", "rem"},
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

	// Get branches for a given repository
	branches, err := util.Spinner("Retrieving branches...", func() ([]string, error) {
		branches, err := ccClient.GetBranches(cfg.Repository)
		if err != nil {
			return nil, err
		}
		return branches, nil
	})
	util.ExitOnErr(err)

	// Prompt for branches to remove
	branchSelections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(branches).Show("Select branches to remove")
	util.ExitOnErr(err)

	// Remove branches
	results, _ := util.Spinner("Removing...", func() ([]client.Result[string], error) {
		results := ccClient.DeleteBranches(cfg.Repository, branchSelections)
		return results, nil
	})

	failCount := 0
	for _, r := range results {
		if r.Err != nil {
			pterm.Error.Printf("Failed to remove branch [%s]: %s\n", r.Result, r.Err)
			failCount++
			continue
		}
		pterm.Success.Printf("Successfully removed branch [%s]\n", r.Result)
	}

	if failCount > 0 {
		os.Exit(1)
	}
}
