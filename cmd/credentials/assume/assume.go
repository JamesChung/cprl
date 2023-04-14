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
	cfg, err := config.NewCredentialsConfig(cmd)
	util.ExitOnErr(err)

	stsClient, err := client.NewSTSClient(cfg.AWSProfile)
	util.ExitOnErr(err)

	if cfg.RoleARN == "" {
		iamClient, err := client.NewIAMClient(cfg.AWSProfile)
		util.ExitOnErr(err)
		roles, err := util.Spinner("Getting roles...", func() ([]types.Role, error) {
			roles, err := iamClient.ListRoles()
			if err != nil {
				return nil, err
			}
			return roles, nil
		})
		util.ExitOnErr(err)

		roleNames := make([]string, 0, 10)
		roleMap := make(map[string]string)
		for _, r := range roles {
			roleNames = append(roleNames, aws.ToString(r.RoleName))
			roleMap[aws.ToString(r.RoleName)] = aws.ToString(r.Arn)
		}

		roleName, err := pterm.DefaultInteractiveSelect.
			WithOptions(roleNames).Show("Select role to assume")
		util.ExitOnErr(err)

		cfg.RoleARN = roleMap[roleName]
	}

	if cfg.SessionName == "" {
		cfg.SessionName, err = pterm.DefaultInteractiveTextInput.
			WithDefaultText("Session name").Show()
		util.ExitOnErr(err)
	}

	creds, err := util.Spinner("Acquiring credentials...", func() (*sts.AssumeRoleOutput, error) {
		creds, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(strings.Trim(cfg.RoleARN, " ")),
			RoleSessionName: aws.String(strings.Trim(cfg.SessionName, " ")),
		})
		if err != nil {
			return nil, err
		}
		return creds, nil
	})
	util.ExitOnErr(err)

	if cfg.OutputProfile == "" {
		cfg.OutputProfile, err = pterm.DefaultInteractiveTextInput.
			WithDefaultText("New AWS profile name").Show()
		util.ExitOnErr(err)
	}

	msg, err := util.Spinner("Writing credentials...", func() (string, error) {
		err = util.WriteCredentials(cfg.OutputProfile, aws.Credentials{
			AccessKeyID:     aws.ToString(creds.Credentials.AccessKeyId),
			SecretAccessKey: aws.ToString(creds.Credentials.SecretAccessKey),
			SessionToken:    aws.ToString(creds.Credentials.SessionToken),
		})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%s] was saved to credentials", cfg.OutputProfile), nil
	})
	util.ExitOnErr(err)
	pterm.Success.Println(msg)
}
