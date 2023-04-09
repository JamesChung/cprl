package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type CloudWatchClient struct {
	Client *cloudwatchlogs.Client
}

func newCloudWatchClient(profile string) (*cloudwatchlogs.Client, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(
		ctx, config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	return cloudwatchlogs.NewFromConfig(cfg), nil
}

func NewCloudWatchClient(profile string) (*CloudWatchClient, error) {
	client, err := newCloudWatchClient(profile)
	if err != nil {
		return nil, err
	}

	c := &CloudWatchClient{
		Client: client,
	}
	return c, nil
}
