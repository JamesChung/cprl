package codecommit

import (
	"fmt"

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

func GetAuthorARNFlag(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "author-arn")
	if err != nil {
		return "", err
	}
	return str, nil
}
