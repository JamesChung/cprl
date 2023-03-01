package codecommit

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JamesChung/cprl/internal/errs"
	"github.com/JamesChung/cprl/pkg/util"
)

func GetRepositories(profile string) ([]string, error) {
	repositories := viper.GetStringSlice(
		fmt.Sprintf(
			"%s.services.codecommit.repositories",
			profile,
		),
	)
	if repositories == nil {
		return nil, errs.MissingRepositoriesErr
	}
	return repositories, nil
}

func GetAuthorARN(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "author-arn")
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetClosed(cmd *cobra.Command) (types.PullRequestStatusEnum, error) {
	closed, err := util.GetFlagBool(cmd, "closed")
	if err != nil {
		return types.PullRequestStatusEnumOpen, err
	}
	if closed {
		return types.PullRequestStatusEnumClosed, nil
	}
	return types.PullRequestStatusEnumOpen, nil
}
