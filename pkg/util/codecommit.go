package util

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/pterm/pterm"
	"golang.org/x/exp/slices"

	"github.com/JamesChung/cprl/internal/errs"
	"github.com/JamesChung/cprl/pkg/client"
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
			if err != nil {
				errs.ErrorCh <- err
				return
			}
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
				if err != nil {
					errs.ErrorCh <- err
					return
				}
				ch <- info
			}
		}(ids)
	}
	return ch
}

func GetPullRequestsSlice(input PullRequestInput) ([]*codecommit.GetPullRequestOutput, error) {
	prs := []*codecommit.GetPullRequestOutput{}
	prCh := GetPullRequests(input)
	for {
		select {
		case pr, ok := <-prCh:
			if !ok {
				return prs, nil
			}
			prs = append(prs, pr)
		case err := <-errs.ErrorCh:
			return nil, err
		}
	}
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

func PRsToTable(headers []string, ch <-chan *codecommit.GetPullRequestOutput) *pterm.TablePrinter {
	data := pterm.TableData{headers}
	for pr := range ch {
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
