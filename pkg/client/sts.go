package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSClient struct {
	Client *sts.Client
}

// newSTSClient creates a localized STS client from the given profile
func newSTSClient(profile string) (*sts.Client, error) {
	cfg, err := GetProfileConfig(profile)
	if err != nil {
		return nil, err
	}
	return sts.NewFromConfig(cfg), nil
}

// NewSTSClient returns a convenient wrapper around the actual AWS STS SDK client
func NewSTSClient(profile string) (*STSClient, error) {
	c, err := newSTSClient(profile)
	if err != nil {
		return nil, err
	}

	stsClient := &STSClient{
		Client: c,
	}
	return stsClient, nil
}

// GetCallerIdentity returns the caller identity of the profile the associated client
// was created with. i.e. from `NewSTSClient`
func (s *STSClient) GetCallerIdentity() (*sts.GetCallerIdentityOutput, error) {
	ctx := context.Background()
	out, err := s.Client.GetCallerIdentity(
		ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *STSClient) AssumeRole(params *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	ctx := context.Background()
	out, err := s.Client.AssumeRole(ctx, params)
	if err != nil {
		return nil, err
	}
	return out, nil
}
