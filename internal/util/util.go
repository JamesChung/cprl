package util

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/spf13/cobra"

	"github.com/JamesChung/cprl/internal/config"
	"github.com/JamesChung/cprl/internal/cprlerrs"
	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

var (
	authorARN string
	cfg       *config.Config
	ccClient  *client.CodeCommitClient
)

func Initialize(cmd *cobra.Command) {
	profile, err := GetProfileFlag(cmd)
	cprlerrs.ExitOnErr(cmd, err)

	cfg, err = config.GetConfig(profile)
	cprlerrs.ExitOnErr(cmd, err)

	awsProfile, err := GetAWSProfileFlag(cmd)
	cprlerrs.ExitOnErr(cmd, err)

	ccClient, err = client.NewCodeCommitClient(awsProfile)
	cprlerrs.ExitOnErr(cmd, err)

	authorARN, err = GetAuthorARNFlag(cmd)
	cprlerrs.ExitOnErr(cmd, err)
}

func GetProfileFlag(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "profile")
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetAWSProfileFlag(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "aws-profile")
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetAuthorARNFlag(cmd *cobra.Command) (string, error) {
	str, err := util.GetFlagString(cmd, "author-arn")
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetPullRequestIDs(cmd *cobra.Command, status types.PullRequestStatusEnum) <-chan []string {
	wg := sync.WaitGroup{}
	ch := make(chan []string, 10)
	defer func() {
		go func() {
			defer close(ch)
			wg.Wait()
		}()
	}()
	for _, repo := range cfg.Repositories {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			ids, err := ccClient.ListPRs(repo, authorARN, status)
			cprlerrs.ExitOnErr(cmd, err)
			ch <- ids
		}(repo)
	}
	return ch
}

func GetPullRequests(cmd *cobra.Command, status types.PullRequestStatusEnum) <-chan *codecommit.GetPullRequestOutput {
	chIDs := GetPullRequestIDs(cmd, status)
	wg := sync.WaitGroup{}
	ch := make(chan *codecommit.GetPullRequestOutput, 10)
	defer func() {
		go func() {
			defer close(ch)
			wg.Wait()
		}()
	}()
	for ids := range chIDs {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				info, err := ccClient.GetPRInfo(id)
				cprlerrs.ExitOnErr(cmd, err)
				ch <- info
			}
		}(ids)
	}
	return ch
}
