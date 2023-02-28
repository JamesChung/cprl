package util

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func AddGroup(parent *cobra.Command, title string, cmds ...*cobra.Command) {
	group := &cobra.Group{
		Title: title,
		ID:    title,
	}
	parent.AddGroup(group)
	for _, cmd := range cmds {
		cmd.GroupID = group.ID
		parent.AddCommand(cmd)
	}
}

func Basename(str string) string {
	s := strings.Split(str, "/")
	return s[len(s)-1]
}

func GetFlagString(cmd *cobra.Command, str string) (string, error) {
	val, err := cmd.Flags().GetString(str)
	if err != nil {
		return "", err
	}
	return val, nil
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
				Basename(aws.ToString(pr.PullRequest.AuthorArn)),
				aws.ToString(pr.PullRequest.Title),
				Basename(aws.ToString(t.SourceReference)),
				Basename(aws.ToString(t.DestinationReference)),
				aws.ToTime(pr.PullRequest.CreationDate).Format(time.DateOnly),
				aws.ToTime(pr.PullRequest.LastActivityDate).Format(time.DateOnly),
			})
		}
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
