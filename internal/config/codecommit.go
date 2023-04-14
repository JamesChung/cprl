package config

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/JamesChung/cprl/pkg/util"
)

type CodeCommitConfig struct {
	*Config
	AuthorARN  string
	Repository string
	PRStatus   types.PullRequestStatusEnum
}

func NewCodeCommitConfig(cmd *cobra.Command) (*CodeCommitConfig, error) {
	c, err := NewConfig(cmd)
	if err != nil {
		return nil, err
	}
	cfg := &CodeCommitConfig{}
	cfg.Config = c
	cfg.AuthorARN, _ = cmd.Flags().GetString("author-arn")
	cfg.Repository, _ = cmd.Flags().GetString("repository")
	cfg.PRStatus = getPRStatus(cmd)
	return cfg, nil
}

func GetRepositories(profile string) ([]string, error) {
	repositories := viper.GetStringSlice(
		fmt.Sprintf(
			"%s.services.codecommit.repositories",
			profile,
		),
	)
	if repositories == nil {
		return nil, errors.New("cprl.yaml is missing [repositories]")
	}
	return repositories, nil
}

func getPRStatus(cmd *cobra.Command) types.PullRequestStatusEnum {
	closed, _ := util.GetFlagBool(cmd, "closed")
	if closed {
		return types.PullRequestStatusEnumClosed
	}
	return types.PullRequestStatusEnumOpen
}
