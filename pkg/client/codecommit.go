package client

import (
	"context"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"

	"github.com/JamesChung/cprl/internal/diff"
)

type CodeCommitClient struct {
	Client *codecommit.Client
}

// PRMap ...
type PRMap map[string]*codecommit.GetPullRequestOutput

// ListPRs returns a list of CodeCommit PR IDs
func (c *CodeCommitClient) ListPRs(repositoryName, authorArn string, status types.PullRequestStatusEnum) ([]string, error) {
	ids := []string{}
	ctx := context.Background()
	p := codecommit.NewListPullRequestsPaginator(
		c.Client, &codecommit.ListPullRequestsInput{
			RepositoryName:    aws.String(repositoryName),
			AuthorArn:         NullableString(authorArn),
			PullRequestStatus: status,
		})
	for p.HasMorePages() {
		o, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		ids = append(ids, o.PullRequestIds...)
	}
	return ids, nil
}

func (c *CodeCommitClient) GetPRInfo(prID string) (*codecommit.GetPullRequestOutput, error) {
	ctx := context.Background()
	pr, err := c.Client.GetPullRequest(
		ctx, &codecommit.GetPullRequestInput{
			PullRequestId: aws.String(prID),
		})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (c *CodeCommitClient) GetBranches(repoName string) ([]string, error) {
	branches := []string{}
	ctx := context.Background()
	p := codecommit.NewListBranchesPaginator(
		c.Client, &codecommit.ListBranchesInput{
			RepositoryName: aws.String(repoName),
		})
	for p.HasMorePages() {
		b, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		branches = append(branches, b.Branches...)
	}
	return branches, nil
}

func (c *CodeCommitClient) DeleteBranch(repo string, branch string) (*codecommit.DeleteBranchOutput, error) {
	ctx := context.Background()
	res, err := c.Client.DeleteBranch(ctx, &codecommit.DeleteBranchInput{
		RepositoryName: aws.String(repo),
		BranchName:     aws.String(branch),
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteBranches ...
func (c *CodeCommitClient) DeleteBranches(repo string, branches []string) []Result[string] {
	resCh := make(chan Result[string], runtime.NumCPU())
	wg := sync.WaitGroup{}
	results := make([]Result[string], 0, 10)
	for _, branch := range branches {
		wg.Add(1)
		go func(branch string) {
			defer wg.Done()
			res, err := c.DeleteBranch(repo, branch)
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

func (c *CodeCommitClient) CreatePR(targets []types.Target, title, desc string) (*codecommit.CreatePullRequestOutput, error) {
	ctx := context.Background()
	p, err := c.Client.CreatePullRequest(
		ctx, &codecommit.CreatePullRequestInput{
			Targets:     targets,
			Title:       aws.String(title),
			Description: NullableString(desc),
		})
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ApprovePRs ...
func (c *CodeCommitClient) ApprovePRs(prMap PRMap, prSelections []string) []Result[string] {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	ch := make(chan Result[string], 10)
	for _, v := range prSelections {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			_, err := c.Client.UpdatePullRequestApprovalState(
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
	results := make([]Result[string], 0, 10)
	for res := range ch {
		results = append(results, res)
	}
	return results
}

// ClosePRs ...
func (c *CodeCommitClient) ClosePRs(prMap PRMap, prSelections []string) []Result[string] {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	ch := make(chan Result[string], 10)
	for _, v := range prSelections {
		wg.Add(1)
		go func(v string) {
			defer wg.Done()
			_, err := c.Client.UpdatePullRequestStatus(
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
	results := make([]Result[string], 0, 10)
	for res := range ch {
		results = append(results, res)
	}
	return results
}

func (c *CodeCommitClient) GetDifferences(repositoryName, beforeCommitSpecifier, afterCommitSpecifier *string) ([]*codecommit.GetDifferencesOutput, error) {
	ctx := context.Background()
	diffs := make([]*codecommit.GetDifferencesOutput, 0)
	p := codecommit.NewGetDifferencesPaginator(c.Client,
		&codecommit.GetDifferencesInput{
			RepositoryName:        repositoryName,
			BeforeCommitSpecifier: beforeCommitSpecifier,
			AfterCommitSpecifier:  afterCommitSpecifier,
		})
	for p.HasMorePages() {
		o, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, o)
	}
	return diffs, nil
}

// PullRequestInput ...
type PullRequestInput struct {
	AuthorARN    string
	Repositories []string
	Status       types.PullRequestStatusEnum
}

// GetPullRequestIDs ...
func (c *CodeCommitClient) GetPullRequestIDs(input PullRequestInput) ([][]string, error) {
	ch := make(chan Result[[]string], 10)
	wg := sync.WaitGroup{}
	for _, repo := range input.Repositories {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			ids, err := c.ListPRs(
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

// GetPullRequestInfoFromIDs ...
func (c *CodeCommitClient) GetPullRequestInfoFromIDs(input [][]string) ([]*codecommit.GetPullRequestOutput, error) {
	ch := make(chan Result[*codecommit.GetPullRequestOutput], 10)
	wg := sync.WaitGroup{}
	for _, ids := range input {
		wg.Add(1)
		go func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				info, err := c.GetPRInfo(id)
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
func (c *CodeCommitClient) GetPullRequests(input PullRequestInput) ([]*codecommit.GetPullRequestOutput, error) {
	ids, err := c.GetPullRequestIDs(input)
	if err != nil {
		return nil, err
	}
	prInfoList, err := c.GetPullRequestInfoFromIDs(ids)
	if err != nil {
		return nil, err
	}
	return prInfoList, err
}

// filterDiffErrors ...
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

// GenerateDiffs ...
func (c *CodeCommitClient) GenerateDiffs(repo string, diffOut []*codecommit.GetDifferencesOutput) []Result[[]byte] {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	ch := make(chan Result[[]byte], runtime.NumCPU())
	for _, do := range diffOut {
		for _, d := range do.Differences {
			wg.Add(1)
			go func(d types.Difference) {
				defer wg.Done()
				sleep := ExponentialBackoff(time.Millisecond*100, time.Second*3)
			retry:
				switch d.ChangeType {
				case types.ChangeTypeEnumModified:
					bob, err := c.Client.GetBlob(ctx, &codecommit.GetBlobInput{
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
					boa, err := c.Client.GetBlob(ctx, &codecommit.GetBlobInput{
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
					boa, err := c.Client.GetBlob(ctx, &codecommit.GetBlobInput{
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
					bob, err := c.Client.GetBlob(ctx, &codecommit.GetBlobInput{
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

func newCodeCommitClient(profile string) (*codecommit.Client, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(
		ctx, config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	return codecommit.NewFromConfig(cfg), nil
}

func NewCodeCommitClient(profile string) (*CodeCommitClient, error) {
	client, err := newCodeCommitClient(profile)
	if err != nil {
		return nil, err
	}

	c := &CodeCommitClient{
		Client: client,
	}
	return c, nil
}
