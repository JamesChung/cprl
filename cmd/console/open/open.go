package open

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/internal/config/services/console"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "opens AWS console"

	example = templates.Examples(`
	Open console with default cprl profile:
	$ cprl console open

	Open console with a specific aws profile:
	$ cprl console open --aws-profile=dev`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.BoolP(
		"prompt",
		"p",
		false,
		"prompt user with suggested profiles",
	)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "open",
		Aliases: []string{"o"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(cmd.PersistentFlags())
	return cmd
}

func runCmd(cmd *cobra.Command, args []string) {
	profile, err := config.GetProfile(cmd)
	util.ExitOnErr(err)

	var awsProfile string
	prompt, err := cmd.Flags().GetBool("prompt")
	util.ExitOnErr(err)
	if prompt {
		profiles, err := util.GetAllProfiles()
		if err != nil {
			util.ExitOnErr(err)
		}
		awsProfile, err = pterm.DefaultInteractiveSelect.
			WithOptions(profiles).Show("Select a profile")
		util.ExitOnErr(err)
	} else {
		awsProfile, err = config.GetAWSProfile(cmd)
		util.ExitOnErr(err)
	}

	isGovCloud, err := console.IsGovCloud(cmd, profile)
	util.ExitOnErr(err)
	creds, err := util.GetCredentials(awsProfile)
	util.ExitOnErr(err)
	loginURL, err := util.GenerateLoginURL(creds, isGovCloud)
	util.ExitOnErr(err)
	err = util.OpenBrowser(loginURL.String())
	util.ExitOnErr(err)
}
