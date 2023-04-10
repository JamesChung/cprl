package util

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/pterm/pterm"
	"golang.org/x/exp/slices"
)

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
