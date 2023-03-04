package util

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/pterm/pterm"
	"golang.org/x/exp/slices"

	"github.com/JamesChung/cprl/pkg/client"
)

type PullRequestInput struct {
	AuthorARN    string
	Client       *client.CodeCommitClient
	Repositories []string
	Status       types.PullRequestStatusEnum
}

func GetPullRequestIDs(input PullRequestInput) ([][]string, error) {
	type result struct {
		IDs   []string
		Error error
	}
	ch := make(chan result, 10)
	wg := sync.WaitGroup{}
	for _, repo := range input.Repositories {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			ids, err := input.Client.ListPRs(
				repo, input.AuthorARN, input.Status,
			)
			if err != nil {
				ch <- result{nil, err}
				return
			}
			ch <- result{ids, nil}
		}(repo)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	response := make([][]string, 0, len(input.Repositories))
	for ids := range ch {
		if ids.Error != nil {
			return nil, ids.Error
		}
		response = append(response, ids.IDs)
	}
	return response, nil
}

func GetPullRequestInfoFromIDs(ccClient *client.CodeCommitClient, input [][]string) ([]*codecommit.GetPullRequestOutput, error) {
	type result struct {
		PRInfo *codecommit.GetPullRequestOutput
		Error  error
	}
	ch := make(chan result, 10)
	wg := sync.WaitGroup{}
	for _, ids := range input {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				info, err := ccClient.GetPRInfo(id)
				if err != nil {
					ch <- result{nil, err}
					return
				}
				ch <- result{info, nil}
			}
		}(ids)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	prList := make([]*codecommit.GetPullRequestOutput, 0)
	for r := range ch {
		if r.Error != nil {
			return nil, r.Error
		}
		prList = append(prList, r.PRInfo)
	}
	return prList, nil
}

// GetPullRequests combines GetPullRequestIDs & GetPullRequestInfoFromIDs into one call
func GetPullRequests(input PullRequestInput) ([]*codecommit.GetPullRequestOutput, error) {
	ids, err := GetPullRequestIDs(input)
	if err != nil {
		return nil, err
	}
	prInfoList, err := GetPullRequestInfoFromIDs(input.Client, ids)
	if err != nil {
		return nil, err
	}
	return prInfoList, err
}

func GenerateTableHeaders(headers []string) []string {
	s := make([]string, 0, len(headers))
	if slices.Contains(headers, "Repository") {
		s = append(s, "Repository")
	}
	if slices.Contains(headers, "Author") {
		s = append(s, "Author")
	}
	if slices.Contains(headers, "ID") {
		s = append(s, "ID")
	}
	if slices.Contains(headers, "Title") {
		s = append(s, "Title")
	}
	if slices.Contains(headers, "Source") {
		s = append(s, "Source")
	}
	if slices.Contains(headers, "Destination") {
		s = append(s, "Destination")
	}
	if slices.Contains(headers, "CreationDate") {
		s = append(s, "CreationDate")
	}
	if slices.Contains(headers, "LastActivityDate") {
		s = append(s, "LastActivityDate")
	}
	return s
}

func PRsToTable(headers []string, prList []*codecommit.GetPullRequestOutput) *pterm.TablePrinter {
	data := pterm.TableData{headers}
	for _, pr := range prList {
		for _, t := range pr.PullRequest.PullRequestTargets {
			row := make([]string, 0, len(headers))
			if slices.Contains(headers, "Repository") {
				row = append(row, aws.ToString(t.RepositoryName))
			}
			if slices.Contains(headers, "Author") {
				row = append(row, Basename(aws.ToString(pr.PullRequest.AuthorArn)))
			}
			if slices.Contains(headers, "ID") {
				row = append(row, Basename(aws.ToString(pr.PullRequest.PullRequestId)))
			}
			if slices.Contains(headers, "Title") {
				row = append(row, aws.ToString(pr.PullRequest.Title))
			}
			if slices.Contains(headers, "Source") {
				row = append(row, Basename(aws.ToString(t.SourceReference)))
			}
			if slices.Contains(headers, "Destination") {
				row = append(row, Basename(aws.ToString(t.DestinationReference)))
			}
			if slices.Contains(headers, "CreationDate") {
				row = append(row, aws.ToTime(pr.PullRequest.CreationDate).Format(time.DateOnly))
			}
			if slices.Contains(headers, "LastActivityDate") {
				row = append(row, aws.ToTime(pr.PullRequest.LastActivityDate).Format(time.DateOnly))
			}
			data = append(data, row)
		}
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data)
}
