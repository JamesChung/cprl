package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
)

func TestGenerateTableHeaders(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "all headers in order",
			input: []string{"Repository", "Author", "ID", "Title", "Source", "Destination", "CreationDate", "LastActivityDate"},
			want:  []string{"Repository", "Author", "ID", "Title", "Source", "Destination", "CreationDate", "LastActivityDate"},
		},
		{
			name:  "subset of headers",
			input: []string{"ID", "Title", "Author"},
			want:  []string{"Author", "ID", "Title"},
		},
		{
			name:  "single header",
			input: []string{"Title"},
			want:  []string{"Title"},
		},
		{
			name:  "headers in wrong order get reordered",
			input: []string{"Title", "Repository", "ID"},
			want:  []string{"Repository", "ID", "Title"},
		},
		{
			name:  "empty input",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "duplicate headers",
			input: []string{"Title", "Title", "ID"},
			want:  []string{"ID", "Title"},
		},
		{
			name:  "invalid headers are ignored",
			input: []string{"InvalidHeader", "Title", "ID"},
			want:  []string{"ID", "Title"},
		},
		{
			name:  "date headers only",
			input: []string{"CreationDate", "LastActivityDate"},
			want:  []string{"CreationDate", "LastActivityDate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTableHeaders(tt.input)

			if len(got) != len(tt.want) {
				t.Errorf("GenerateTableHeaders() returned %d headers, want %d", len(got), len(tt.want))
				t.Errorf("got: %v", got)
				t.Errorf("want: %v", tt.want)
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GenerateTableHeaders()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestPRsToTable(t *testing.T) {
	// Helper function to create a test PR
	createTestPR := func(repo, author, id, title, source, dest string, created, lastActivity time.Time) *codecommit.GetPullRequestOutput {
		return &codecommit.GetPullRequestOutput{
			PullRequest: &types.PullRequest{
				PullRequestId:    aws.String(id),
				Title:            aws.String(title),
				AuthorArn:        aws.String("arn:aws:iam::123456789012:user/" + author),
				CreationDate:     aws.Time(created),
				LastActivityDate: aws.Time(lastActivity),
				PullRequestTargets: []types.PullRequestTarget{
					{
						RepositoryName:       aws.String(repo),
						SourceReference:      aws.String("refs/heads/" + source),
						DestinationReference: aws.String("refs/heads/" + dest),
					},
				},
			},
		}
	}

	createdTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	activityTime := time.Date(2024, 1, 16, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name       string
		headers    []string
		prList     []*codecommit.GetPullRequestOutput
		wantRows   int
		wantCols   int
		checkCells map[string]string // map of "row,col" -> expected value
	}{
		{
			name:    "single PR with all headers",
			headers: []string{"Repository", "Author", "ID", "Title", "Source", "Destination", "CreationDate", "LastActivityDate"},
			prList: []*codecommit.GetPullRequestOutput{
				createTestPR("my-repo", "john", "123", "Fix bug", "feature-branch", "main", createdTime, activityTime),
			},
			wantRows: 2, // header + 1 data row
			wantCols: 8,
			checkCells: map[string]string{
				"1,0": "my-repo",
				"1,1": "john",
				"1,2": "123",
				"1,3": "Fix bug",
				"1,4": "feature-branch",
				"1,5": "main",
				"1,6": "2024-01-15",
				"1,7": "2024-01-16",
			},
		},
		{
			name:    "subset of headers",
			headers: []string{"ID", "Title", "Author"},
			prList: []*codecommit.GetPullRequestOutput{
				createTestPR("my-repo", "alice", "456", "Add feature", "dev", "main", createdTime, activityTime),
			},
			wantRows: 2,
			wantCols: 3,
			checkCells: map[string]string{
				"1,0": "alice",
				"1,1": "456",
				"1,2": "Add feature",
			},
		},
		{
			name:    "multiple PRs",
			headers: []string{"ID", "Title"},
			prList: []*codecommit.GetPullRequestOutput{
				createTestPR("repo1", "user1", "100", "First PR", "br1", "main", createdTime, activityTime),
				createTestPR("repo2", "user2", "200", "Second PR", "br2", "main", createdTime, activityTime),
			},
			wantRows: 3, // header + 2 data rows
			wantCols: 2,
			checkCells: map[string]string{
				"1,0": "100",
				"1,1": "First PR",
				"2,0": "200",
				"2,1": "Second PR",
			},
		},
		{
			name:     "empty PR list",
			headers:  []string{"ID", "Title"},
			prList:   []*codecommit.GetPullRequestOutput{},
			wantRows: 1, // just header
			wantCols: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := PRsToTable(tt.headers, tt.prList)

			// Get the table data
			data := table.Data

			if len(data) != tt.wantRows {
				t.Errorf("PRsToTable() returned %d rows, want %d", len(data), tt.wantRows)
				return
			}

			if len(data) > 0 && len(data[0]) != tt.wantCols {
				t.Errorf("PRsToTable() returned %d columns, want %d", len(data[0]), tt.wantCols)
				return
			}

			// Check specific cells if provided
			for cellKey, expectedValue := range tt.checkCells {
				var row, col int
				_, err := fmt.Sscanf(cellKey, "%d,%d", &row, &col)
				if err != nil {
					t.Fatalf("Invalid cell key format: %s", cellKey)
				}

				if row >= len(data) {
					t.Errorf("Row %d out of bounds (total rows: %d)", row, len(data))
					continue
				}

				if col >= len(data[row]) {
					t.Errorf("Column %d out of bounds for row %d (total cols: %d)", col, row, len(data[row]))
					continue
				}

				got := data[row][col]
				if got != expectedValue {
					t.Errorf("Cell [%d,%d] = %q, want %q", row, col, got, expectedValue)
				}
			}
		})
	}
}

func TestPRsToTable_MultipleTargets(t *testing.T) {
	// Test PR with multiple targets (each target should create a row)
	pr := &codecommit.GetPullRequestOutput{
		PullRequest: &types.PullRequest{
			PullRequestId:    aws.String("999"),
			Title:            aws.String("Multi-target PR"),
			AuthorArn:        aws.String("arn:aws:iam::123456789012:user/bob"),
			CreationDate:     aws.Time(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
			LastActivityDate: aws.Time(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
			PullRequestTargets: []types.PullRequestTarget{
				{
					RepositoryName:       aws.String("repo1"),
					SourceReference:      aws.String("refs/heads/feature"),
					DestinationReference: aws.String("refs/heads/main"),
				},
				{
					RepositoryName:       aws.String("repo2"),
					SourceReference:      aws.String("refs/heads/feature"),
					DestinationReference: aws.String("refs/heads/develop"),
				},
			},
		},
	}

	headers := []string{"Repository", "ID", "Title"}
	table := PRsToTable(headers, []*codecommit.GetPullRequestOutput{pr})

	// Should have header + 2 data rows (one per target)
	if len(table.Data) != 3 {
		t.Errorf("PRsToTable() with 2 targets returned %d rows, want 3", len(table.Data))
	}

	// Both rows should have the same ID and Title
	if len(table.Data) >= 3 {
		if table.Data[1][1] != "999" || table.Data[2][1] != "999" {
			t.Errorf("ID not consistent across target rows")
		}
		if table.Data[1][2] != "Multi-target PR" || table.Data[2][2] != "Multi-target PR" {
			t.Errorf("Title not consistent across target rows")
		}
		// But different repositories
		if table.Data[1][0] != "repo1" {
			t.Errorf("First target repository = %q, want %q", table.Data[1][0], "repo1")
		}
		if table.Data[2][0] != "repo2" {
			t.Errorf("Second target repository = %q, want %q", table.Data[2][0], "repo2")
		}
	}
}

func TestBasename_InCodeCommitContext(t *testing.T) {
	// Test Basename with CodeCommit-specific strings
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "extract username from ARN",
			input: "arn:aws:iam::123456789012:user/john.doe",
			want:  "john.doe",
		},
		{
			name:  "extract PR ID",
			input: "arn:aws:codecommit:us-east-1:123456789012:pr/42",
			want:  "42",
		},
		{
			name:  "extract branch from refs",
			input: "refs/heads/feature-branch",
			want:  "feature-branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Basename(tt.input)
			if got != tt.want {
				t.Errorf("Basename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
