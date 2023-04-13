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
		"interactive",
		"i",
		false,
		"interactive prompt with suggested profiles",
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

type openConfig struct {
	config.Config
	Interactive bool
	IsGovCloud  bool
}

func newOpenConfig(cmd *cobra.Command) (*openConfig, error) {
	c, err := config.NewConfig(cmd)
	if err != nil {
		return nil, err
	}
	cfg := &openConfig{}
	cfg.Config = *c
	interactive, err := cmd.Flags().GetBool("interactive")
	if err != nil {
		return nil, err
	}
	cfg.Interactive = interactive
	isGovCloud, err := console.IsGovCloud(cmd, cfg.Profile)
	if err != nil {
		return nil, err
	}
	cfg.IsGovCloud = isGovCloud
	return cfg, nil
}

func runCmd(cmd *cobra.Command, args []string) {
	cfg, err := newOpenConfig(cmd)
	util.ExitOnErr(err)

	if cfg.Interactive {
		profiles, err := util.GetAllProfiles()
		if err != nil {
			util.ExitOnErr(err)
		}
		cfg.AWSProfile, err = pterm.DefaultInteractiveSelect.
			WithOptions(profiles).Show("Select a profile")
		util.ExitOnErr(err)
	}

	creds, err := util.GetCredentials(cfg.AWSProfile)
	util.ExitOnErr(err)
	loginURL, err := util.GenerateLoginURL(creds, cfg.IsGovCloud)
	util.ExitOnErr(err)
	err = util.OpenBrowser(loginURL.String())
	util.ExitOnErr(err)
}
