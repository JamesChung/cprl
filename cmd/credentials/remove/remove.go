package remove

import (
	"github.com/JamesChung/cprl/pkg/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	shortMessage = "remove AWS credentials"

	example = templates.Examples(`
	$ cprl credentials remove
	Select profiles to remove:
	> [✗] dev
	enter: select | tab: confirm | left: none | right: all| type to filter`)
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "rm", "rem"},
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
		WithOptions(profiles).Show("Select profiles to remove")
	util.ExitOnErr(err)
	err = util.RemoveProfiles(selections)
	util.ExitOnErr(err)
}
