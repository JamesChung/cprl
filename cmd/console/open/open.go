package open

import (
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/internal/config/services/console"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "opens AWS console"

	example = templates.Examples(`
	cprl console open
	cprl console open --aws-profile=dev-account`)
)

func NewCmd() *cobra.Command {
	consoleCmd := &cobra.Command{
		Use:     "open",
		Aliases: []string{"o"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	return consoleCmd
}

func runCmd(cmd *cobra.Command, args []string) {
	profile, err := config.GetProfile(cmd)
	util.ExitOnErr(err)
	awsProfile, err := config.GetAWSProfile(cmd)
	util.ExitOnErr(err)
	isGovCloud, err := console.IsGovCloud(cmd, profile)
	util.ExitOnErr(err)
	cmd.Println(isGovCloud, awsProfile)
}
