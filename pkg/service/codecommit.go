package service

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/pterm/pterm"

	"github.com/JamesChung/cprl/pkg/client"
	"github.com/JamesChung/cprl/pkg/util"
)

type PullRequestInput struct {
	AuthorARN    string
	Client       *client.CodeCommitClient
	Repositories []string
	Status       types.PullRequestStatusEnum
}

func GetPullRequestIDs(input PullRequestInput) <-chan []string {
	wg := sync.WaitGroup{}
	ch := make(chan []string, 10)
	defer func() {
		go func() {
			defer close(ch)
			wg.Wait()
		}()
	}()
	for _, repo := range input.Repositories {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			ids, err := input.Client.ListPRs(
				repo, input.AuthorARN, input.Status,
			)
			util.ExitOnErr(err)
			ch <- ids
		}(repo)
	}
	return ch
}

func GetPullRequests(input PullRequestInput) <-chan *codecommit.GetPullRequestOutput {
	chIDs := GetPullRequestIDs(input)
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
				info, err := input.Client.GetPRInfo(id)
				util.ExitOnErr(err)
				ch <- info
			}
		}(ids)
	}
	return ch
}

func PRsToTable(ch <-chan *codecommit.GetPullRequestOutput) {
	data := pterm.TableData{{
		"Repository",
		"Author",
		"Title",
		"Source",
		"Destination",
		"CreationDate",
		"LastActivityDate",
	}}
	for pr := range ch {
		for _, t := range pr.PullRequest.PullRequestTargets {
			data = append(data, []string{
				aws.ToString(t.RepositoryName),
				util.Basename(aws.ToString(pr.PullRequest.AuthorArn)),
				aws.ToString(pr.PullRequest.Title),
				util.Basename(aws.ToString(t.SourceReference)),
				util.Basename(aws.ToString(t.DestinationReference)),
				aws.ToTime(pr.PullRequest.CreationDate).Format(time.DateOnly),
				aws.ToTime(pr.PullRequest.LastActivityDate).Format(time.DateOnly),
			})
		}
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
