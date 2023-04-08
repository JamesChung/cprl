package assume

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
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
	awsProfile, err := config.GetAWSProfile(cmd)
	util.ExitOnErr(err)
	stsClient, err := client.NewSTSClient(awsProfile)
	util.ExitOnErr(err)
	arn, err := pterm.DefaultInteractiveTextInput.WithDefaultText("Role ARN").Show()
	util.ExitOnErr(err)
	sessionName, err := pterm.DefaultInteractiveTextInput.WithDefaultText("Session Name").Show()
	util.ExitOnErr(err)
	creds, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(strings.Trim(arn, " ")),
		RoleSessionName: aws.String(strings.Trim(sessionName, " ")),
	})
	util.ExitOnErr(err)
	profile, err := pterm.DefaultInteractiveTextInput.WithDefaultText("Profile Name").Show()
	util.ExitOnErr(err)
	err = util.WriteCredentials(profile, aws.Credentials{
		AccessKeyID:     aws.ToString(creds.Credentials.AccessKeyId),
		SecretAccessKey: aws.ToString(creds.Credentials.SecretAccessKey),
		SessionToken:    aws.ToString(creds.Credentials.SessionToken),
	})
	util.ExitOnErr(err)
}
