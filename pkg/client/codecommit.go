package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
)

type CodeCommitClient struct {
	Client *codecommit.Client
}

func nilOrString(str string) *string {
	if str == "" {
		return nil
	}
	return aws.String(str)
}

// ListPRs returns a list of CodeCommit PR IDs
func (l *CodeCommitClient) ListPRs(repositoryName, authorArn string, status types.PullRequestStatusEnum) ([]string, error) {
	ids := []string{}
	p := codecommit.NewListPullRequestsPaginator(
		l.Client, &codecommit.ListPullRequestsInput{
			RepositoryName:    aws.String(repositoryName),
			AuthorArn:         nilOrString(authorArn),
			PullRequestStatus: status,
		})
	for p.HasMorePages() {
		o, err := p.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		ids = append(ids, o.PullRequestIds...)
	}
	return ids, nil
}

func (l *CodeCommitClient) GetPRInfo(prID string) (*codecommit.GetPullRequestOutput, error) {
	pr, err := l.Client.GetPullRequest(
		context.Background(),
		&codecommit.GetPullRequestInput{
			PullRequestId: aws.String(prID),
		})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (l *CodeCommitClient) GetBranches(repoName string) ([]string, error) {
	branches := []string{}
	p := codecommit.NewListBranchesPaginator(
		l.Client, &codecommit.ListBranchesInput{
			RepositoryName: aws.String(repoName),
		})
	for p.HasMorePages() {
		b, err := p.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		branches = append(branches, b.Branches...)
	}
	return branches, nil
}

func (l *CodeCommitClient) CreatePR(targets []types.Target, title, desc string) (*codecommit.CreatePullRequestOutput, error) {
	p, err := l.Client.CreatePullRequest(
		context.Background(), &codecommit.CreatePullRequestInput{
			Targets:     targets,
			Title:       aws.String(title),
			Description: nilOrString(desc),
		})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func newCodeCommitClient(profile string) (*codecommit.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithSharedConfigProfile(profile),
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

	l := &CodeCommitClient{
		Client: client,
	}

	return l, nil
}
