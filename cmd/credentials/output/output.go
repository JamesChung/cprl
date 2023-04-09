package output

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "outputs AWS credentials"

	example = templates.Examples(`
	Basic output:
	$ cprl credentials output
	export AWS_ACCESS_KEY_ID=<access key id value>
	export AWS_SECRET_ACCESS_KEY=<secret access key value>
	export AWS_SESSION_TOKEN=<session token value>

	JSON output:
	$ cprl credentials output --json
	{"AccessKeyID":"<access key id value>","SecretAccessKey":...}

	Source credentials as your current session:
	$ source <(cprl credentials output --aws-profile=dev)
	$ aws sts get-caller-identity
	{
		"UserId": "TAG0YY70NST6IUO5KA5XB:cprl",
		"Account": "010203040506",
		"Arn": "arn:aws:sts::010203040506:assumed-role/dev/cprl"
	}`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.Bool(
		"json",
		false,
		"output in json format",
	)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "output",
		Aliases: []string{"o", "ou", "out"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(cmd.PersistentFlags())
	return cmd
}

func runCmd(cmd *cobra.Command, args []string) {
	profiles, err := util.GetAllProfiles()
	util.ExitOnErr(err)

	awsProfile, err := cmd.Flags().GetString("aws-profile")
	util.ExitOnErr(err)

	if awsProfile == "" {
		awsProfile, err = pterm.DefaultInteractiveSelect.
			WithOptions(profiles).Show("Select a profile")
		util.ExitOnErr(err)
	}

	creds, err := util.GetCredentials(awsProfile)
	util.ExitOnErr(err)

	jsonOut, err := cmd.Flags().GetBool("json")
	util.ExitOnErr(err)
	if jsonOut {
		strCreds, err := util.StringifyCredentials(creds)
		util.ExitOnErr(err)
		fmt.Println(strCreds)
		return
	}

	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", creds.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", creds.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", creds.SessionToken)
}