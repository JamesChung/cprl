package pr

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	cc "github.com/JamesChung/cprl/internal/config/services/codecommit"
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
	profile, _ := config.GetProfileFlag(cmd)
	repos, _ := cc.GetRepositories(profile)
	result, err := pterm.DefaultInteractiveSelect.WithOptions(
		repos,
	).Show()
	if err != nil {
		cmd.PrintErrln(err)
	}
	cmd.Println(result)
}
