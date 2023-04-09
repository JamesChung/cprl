package clear

import (
	"github.com/JamesChung/cprl/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortMessage = "clears AWS credentials"

	example = templates.Examples(`cprl credentials clear`)
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clear",
		Aliases: []string{"c", "cl", "clr"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	return cmd
}

func runCmd(cmd *cobra.Command, args []string) {
	profiles, err := util.GetCredentialsProfiles()
	util.ExitOnErr(err)
	selections, err := pterm.DefaultInteractiveMultiselect.
		WithOptions(profiles).Show()
	util.ExitOnErr(err)
	for _, p := range selections {
		err = util.ClearProfile(p)
		if err != nil {
			util.ExitOnErr(err)
		}
	}
}
