package assume

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "assume AWS role"

	example = templates.Examples(`
	Assume role:
	$ cprl credentials assume
	...

	Assume role bypassing input prompts via flags:
	$ cprl --aws-profile=main credentials assume \
		--role-arn=arn:aws:iam::010203040506:role/dev \
		--session-name=cprl --output-profile=dev`)
)

func setPersistentFlags(flags *pflag.FlagSet) {
	flags.String(
		"role-arn",
		"",
		"role ARN of the assuming role",
	)
	flags.String(
		"session-name",
		"",
		"name of the session",
	)
	flags.String(
		"output-profile",
		"",
		"new profile name of the assuming role",
	)
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assume",
		Aliases: []string{"a", "as"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	setPersistentFlags(cmd.PersistentFlags())
	return cmd
}

func runCmd(cmd *cobra.Command, args []string) {
	awsProfile, err := config.GetAWSProfile(cmd)
	util.ExitOnErr(err)

	stsClient, err := client.NewSTSClient(awsProfile)
	util.ExitOnErr(err)

	var roleARN, sessionName, outputProfile string

	roleARN, _ = cmd.LocalFlags().GetString("role-arn")
	if roleARN == "" {
		var iamClient *client.IAMClient
		var roles []types.Role
		util.Spinner("Getting roles...", func() {
			iamClient, err = client.NewIAMClient(awsProfile)
			util.ExitOnErr(err)

			roles, err = iamClient.ListRoles()
			util.ExitOnErr(err)
		})

		roleNames := make([]string, 0, 10)
		roleMap := make(map[string]string)
		for _, r := range roles {
			roleNames = append(roleNames, aws.ToString(r.RoleName))
			roleMap[aws.ToString(r.RoleName)] = aws.ToString(r.Arn)
		}

		roleName, err := pterm.DefaultInteractiveSelect.
			WithOptions(roleNames).Show("Select role to assume")
		util.ExitOnErr(err)

		roleARN = roleMap[roleName]
	}

	sessionName, _ = cmd.LocalFlags().GetString("session-name")
	if sessionName == "" {
		sessionName, err = pterm.DefaultInteractiveTextInput.
			WithDefaultText("Session name").Show()
		util.ExitOnErr(err)
	}

	var creds *sts.AssumeRoleOutput
	util.Spinner("Acquiring credentials...", func() {
		creds, err = stsClient.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(strings.Trim(roleARN, " ")),
			RoleSessionName: aws.String(strings.Trim(sessionName, " ")),
		})
		util.ExitOnErr(err)
	})

	outputProfile, _ = cmd.LocalFlags().GetString("output-profile")
	if outputProfile == "" {
		outputProfile, err = pterm.DefaultInteractiveTextInput.
			WithDefaultText("New AWS profile name").Show()
		util.ExitOnErr(err)
	}

	util.SpinnerWithStatusMsg("Writing credentials...", func() (string, error) {
		err = util.WriteCredentials(outputProfile, aws.Credentials{
			AccessKeyID:     aws.ToString(creds.Credentials.AccessKeyId),
			SecretAccessKey: aws.ToString(creds.Credentials.SecretAccessKey),
			SessionToken:    aws.ToString(creds.Credentials.SessionToken),
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%s] was saved to credentials", outputProfile), nil
	})
}
