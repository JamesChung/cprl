package assume

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortMessage = "assume AWS role"

	example = templates.Examples(`
	cprl credentials assume
	cprl credentials assume --aws-profile=dev-account`)
)

func NewCmd() *cobra.Command {
	consoleCmd := &cobra.Command{
		Use:     "assume",
		Aliases: []string{"a"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	return consoleCmd
}

func runCmd(cmd *cobra.Command, args []string) {
	// profile, err := config.GetProfile(cmd)
	// util.ExitOnErr(err)
	// awsProfile, err := config.GetAWSProfile(cmd)
	// util.ExitOnErr(err)
	// isGovCloud, err := console.IsGovCloud(cmd, profile)
	// util.ExitOnErr(err)
	cmd.Println("Work in Progress...")
}
