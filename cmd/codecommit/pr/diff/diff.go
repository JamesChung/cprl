package diff

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	cc "github.com/JamesChung/cprl/internal/config/services/codecommit"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "Diff PRs"

	longMessage = templates.LongDesc(`
	Diff a PR`)

	example = templates.Examples(`
	cprl codecommit pr diff
	`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	// TODO
}

func NewCmdDiffPR() *cobra.Command {
	diff := &cobra.Command{
		Use:     "diff",
		Aliases: []string{"d"},
		Short:   shortMessage,
		Long:    longMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(diff.PersistentFlags())
	return diff
}

func runCmd(cmd *cobra.Command, args []string) {
	profile, err := config.GetProfile(cmd)
	util.ExitOnErr(err)
	// repos, err := cc.GetRepositories(profile)
	_, err = cc.GetRepositories(profile)
	util.ExitOnErr(err)
	awsProfile, err := config.GetAWSProfile(cmd)
	util.ExitOnErr(err)
	// ccClient, err := client.NewCodeCommitClient(awsProfile)
	_, err = client.NewCodeCommitClient(awsProfile)
	util.ExitOnErr(err)

	cmd.Println("This command is a work in progress...")
}
