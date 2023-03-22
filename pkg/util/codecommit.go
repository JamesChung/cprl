package util

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/pterm/pterm"
	"golang.org/x/exp/slices"

	"github.com/JamesChung/cprl/internal/diff"
	"github.com/JamesChung/cprl/pkg/client"
)

type PullRequestInput struct {
	AuthorARN    string
	Client       *client.CodeCommitClient
	Repositories []string
	Status       types.PullRequestStatusEnum
}

func GetPullRequestIDs(input PullRequestInput) ([][]string, error) {
	ch := make(chan Result[[]string], 10)
	wg := sync.WaitGroup{}
	for _, repo := range input.Repositories {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			ids, err := input.Client.ListPRs(
				repo, input.AuthorARN, input.Status,
			)
			if err != nil {
				ch <- Result[[]string]{nil, err}
				return
			}
			ch <- Result[[]string]{ids, nil}
		}(repo)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	response := make([][]string, 0, len(input.Repositories))
	for ids := range ch {
		if ids.Err != nil {
			return nil, ids.Err
		}
		response = append(response, ids.Result)
	}
	return response, nil
}

func GetPullRequestInfoFromIDs(ccClient *client.CodeCommitClient, input [][]string) ([]*codecommit.GetPullRequestOutput, error) {
	ch := make(chan Result[*codecommit.GetPullRequestOutput], 10)
	wg := sync.WaitGroup{}
	for _, ids := range input {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				info, err := ccClient.GetPRInfo(id)
				if err != nil {
					ch <- Result[*codecommit.GetPullRequestOutput]{nil, err}
					return
				}
				ch <- Result[*codecommit.GetPullRequestOutput]{info, nil}
			}
		}(ids)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	prList := make([]*codecommit.GetPullRequestOutput, 0)
	for r := range ch {
		if r.Err != nil {
			return nil, r.Err
		}
		prList = append(prList, r.Result)
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
				row = append(row, aws.ToTime(pr.PullRequest.CreationDate).Format("2006-01-02"))
			}
			if slices.Contains(headers, "LastActivityDate") {
				row = append(row, aws.ToTime(pr.PullRequest.LastActivityDate).Format("2006-01-02"))
			}
			data = append(data, row)
		}
	}

	return pterm.DefaultTable.WithHasHeader().WithData(data)
}

type PRMap map[string]*codecommit.GetPullRequestOutput

func ApprovePRs(ccClient *client.CodeCommitClient, prMap PRMap, prSelections []string) []Result[string] {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	ch := make(chan Result[string], 10)
	for _, v := range prSelections {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			_, err := ccClient.Client.UpdatePullRequestApprovalState(
				ctx, &codecommit.UpdatePullRequestApprovalStateInput{
					ApprovalState: types.ApprovalStateApprove,
					PullRequestId: prMap[v].PullRequest.PullRequestId,
					RevisionId:    prMap[v].PullRequest.RevisionId,
				})
			if err != nil {
				ch <- Result[string]{v, err}
				return
			}
			ch <- Result[string]{v, nil}
		}(v)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	results := make([]Result[string], 0)
	for res := range ch {
		results = append(results, res)
	}
	return results
}

func ClosePRs(ccClient *client.CodeCommitClient, prMap PRMap, prSelections []string) []Result[string] {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	ch := make(chan Result[string], 10)
	for _, v := range prSelections {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			_, err := ccClient.Client.UpdatePullRequestStatus(
				ctx, &codecommit.UpdatePullRequestStatusInput{
					PullRequestId:     prMap[v].PullRequest.PullRequestId,
					PullRequestStatus: types.PullRequestStatusEnumClosed,
				})
			if err != nil {
				ch <- Result[string]{v, err}
				return
			}
			ch <- Result[string]{v, nil}
		}(v)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	results := make([]Result[string], 0)
	for res := range ch {
		results = append(results, res)
	}
	return results
}

func DeleteBranches(ccClient *client.CodeCommitClient, repo string, branches []string) []Result[string] {
	resCh := make(chan Result[string], runtime.NumCPU())
	wg := sync.WaitGroup{}
	results := make([]Result[string], 0, 10)
	for _, branch := range branches {
		wg.Add(1)
		go func(branch string) {
			defer wg.Done()
			res, err := ccClient.DeleteBranch(repo, branch)
			if err != nil {
				resCh <- Result[string]{branch, err}
				return
			}
			resCh <- Result[string]{aws.ToString(res.DeletedBranch.BranchName), nil}
		}(branch)
	}
	go func() {
		defer close(resCh)
		wg.Wait()
	}()
	for r := range resCh {
		results = append(results, r)
	}
	return results
}

func filterDiffErrors(err error) error {
	switch {
	case strings.Contains(err.Error(), "TLS handshake timeout"):
		return nil
	case strings.Contains(err.Error(), "deserialization failed"):
		return nil
	case strings.Contains(err.Error(), "retry quota exceeded"):
		return nil
	}
	return err
}

func GenerateDiffs(ccClient *client.CodeCommitClient, repo string, diffOut []*codecommit.GetDifferencesOutput) []Result[[]byte] {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	ch := make(chan Result[[]byte], runtime.NumCPU())
	for _, do := range diffOut {
		for _, d := range do.Differences {
			wg.Add(1)
			go func(d types.Difference) {
				defer wg.Done()
				sleep := ExponentialBackoff(time.Millisecond*100, time.Second*2)
			retry:
				switch d.ChangeType {
				case types.ChangeTypeEnumModified:
					bob, err := ccClient.Client.GetBlob(ctx, &codecommit.GetBlobInput{
						BlobId:         d.BeforeBlob.BlobId,
						RepositoryName: aws.String(repo),
					})
					if err != nil {
						if filterDiffErrors(err) == nil {
							sleep()
							goto retry
						}
						ch <- Result[[]byte]{[]byte(aws.ToString(d.BeforeBlob.Path)), err}
						return
					}
					boa, err := ccClient.Client.GetBlob(ctx, &codecommit.GetBlobInput{
						BlobId:         d.AfterBlob.BlobId,
						RepositoryName: aws.String(repo),
					})
					if err != nil {
						if filterDiffErrors(err) == nil {
							sleep()
							goto retry
						}
						ch <- Result[[]byte]{[]byte(aws.ToString(d.AfterBlob.Path)), err}
						return
					}
					diffResult := diff.Diff(
						aws.ToString(d.BeforeBlob.Path),
						bob.Content,
						aws.ToString(d.AfterBlob.Path),
						boa.Content,
					)
					ch <- Result[[]byte]{diffResult, nil}
				case types.ChangeTypeEnumAdded:
					boa, err := ccClient.Client.GetBlob(ctx, &codecommit.GetBlobInput{
						BlobId:         d.AfterBlob.BlobId,
						RepositoryName: aws.String(repo),
					})
					if err != nil {
						if filterDiffErrors(err) == nil {
							sleep()
							goto retry
						}
						ch <- Result[[]byte]{[]byte(aws.ToString(d.AfterBlob.Path)), err}
						return
					}
					diffResult := diff.Diff(
						aws.ToString(d.AfterBlob.Path),
						[]byte{},
						aws.ToString(d.AfterBlob.Path),
						boa.Content,
					)
					ch <- Result[[]byte]{diffResult, nil}
				case types.ChangeTypeEnumDeleted:
					bob, err := ccClient.Client.GetBlob(ctx, &codecommit.GetBlobInput{
						BlobId:         d.BeforeBlob.BlobId,
						RepositoryName: aws.String(repo),
					})
					if err != nil {
						if filterDiffErrors(err) == nil {
							sleep()
							goto retry
						}
						ch <- Result[[]byte]{[]byte(aws.ToString(d.BeforeBlob.Path)), err}
						return
					}
					diffResult := diff.Diff(
						aws.ToString(d.BeforeBlob.Path),
						bob.Content,
						aws.ToString(d.BeforeBlob.Path),
						[]byte{},
					)
					ch <- Result[[]byte]{diffResult, nil}
				}
			}(d)
		}
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	// Poll results
	results := make([]Result[[]byte], 0, 10)
	for res := range ch {
		results = append(results, res)
	}
	return results
}
