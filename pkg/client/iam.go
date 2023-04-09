package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAMClient struct {
	Client *iam.Client
}

func (i *IAMClient) ListRoles() ([]types.Role, error) {
	ctx := context.Background()
	roles := make([]types.Role, 0, 10)
	p := iam.NewListRolesPaginator(
		i.Client, &iam.ListRolesInput{
			PathPrefix: nil,
		})
	for p.HasMorePages() {
		o, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		roles = append(roles, o.Roles...)
	}
	return roles, nil
}

func newIAMClient(profile string) (*iam.Client, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(
		ctx, config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	return iam.NewFromConfig(cfg), nil
}

func NewIAMClient(profile string) (*IAMClient, error) {
	client, err := newIAMClient(profile)
	if err != nil {
		return nil, err
	}

	c := &IAMClient{
		Client: client,
	}
	return c, nil
}
