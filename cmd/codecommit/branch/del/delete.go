package del

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
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
	del := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(del.PersistentFlags())
	return del
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
	var results []*codecommit.DeleteBranchOutput
	var errs []util.Result[string]
	util.Spinner("Deleting...", func() {
		results, errs = util.DeleteBranches(ccClient, repo, branchSelections)
	})

	for _, r := range results {
		pterm.Success.Printf("Successfully deleted branch [%s]\n", aws.ToString(r.DeletedBranch.BranchName))
	}

	if len(errs) > 0 {
		for _, e := range errs {
			pterm.Error.Printf("Failed to delete branch [%s]\n", e.Result)
		}
		util.ExitOnErr(fmt.Errorf("%d branches failed to be deleted\n", len(errs)))
	}
}
