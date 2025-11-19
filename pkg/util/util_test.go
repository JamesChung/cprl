package util

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestBasename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple filename",
			input: "file.txt",
			want:  "file.txt",
		},
		{
			name:  "unix path",
			input: "/path/to/file.txt",
			want:  "file.txt",
		},
		{
			name:  "nested unix path",
			input: "/var/lib/app/data/config.yaml",
			want:  "config.yaml",
		},
		{
			name:  "single directory",
			input: "/home",
			want:  "home",
		},
		{
			name:  "trailing slash",
			input: "/path/to/dir/",
			want:  "",
		},
		{
			name:  "ARN-like string",
			input: "arn:aws:iam::123456789012:user/john",
			want:  "john",
		},
		{
			name:  "multiple slashes",
			input: "path//to///file",
			want:  "file",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "no slashes",
			input: "filename",
			want:  "filename",
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

func TestGetFlagString(t *testing.T) {
	tests := []struct {
		name      string
		flagName  string
		flagValue string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "existing flag with value",
			flagName:  "test-flag",
			flagValue: "test-value",
			wantValue: "test-value",
			wantErr:   false,
		},
		{
			name:      "existing flag with empty value",
			flagName:  "empty-flag",
			flagValue: "",
			wantValue: "",
			wantErr:   false,
		},
		{
			name:      "non-existent flag",
			flagName:  "missing-flag",
			flagValue: "",
			wantValue: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Only add the flag if it's not the "missing-flag" test
			if tt.flagName != "missing-flag" {
				cmd.Flags().String(tt.flagName, "", "test flag")
				cmd.Flags().Set(tt.flagName, tt.flagValue)
			}

			got, err := GetFlagString(cmd, tt.flagName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetFlagString() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("GetFlagString() unexpected error = %v", err)
				}
				if got != tt.wantValue {
					t.Errorf("GetFlagString() = %q, want %q", got, tt.wantValue)
				}
			}
		})
	}
}

func TestGetFlagBool(t *testing.T) {
	tests := []struct {
		name      string
		flagName  string
		flagValue bool
		wantValue bool
		wantErr   bool
	}{
		{
			name:      "existing flag with true",
			flagName:  "bool-flag",
			flagValue: true,
			wantValue: true,
			wantErr:   false,
		},
		{
			name:      "existing flag with false",
			flagName:  "bool-flag",
			flagValue: false,
			wantValue: false,
			wantErr:   false,
		},
		{
			name:      "non-existent flag",
			flagName:  "missing-flag",
			flagValue: false,
			wantValue: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Only add the flag if it's not the "missing-flag" test
			if tt.flagName != "missing-flag" {
				cmd.Flags().Bool(tt.flagName, false, "test flag")
				if tt.flagValue {
					cmd.Flags().Set(tt.flagName, "true")
				}
			}

			got, err := GetFlagBool(cmd, tt.flagName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetFlagBool() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("GetFlagBool() unexpected error = %v", err)
				}
				if got != tt.wantValue {
					t.Errorf("GetFlagBool() = %v, want %v", got, tt.wantValue)
				}
			}
		})
	}
}

func TestSpinner(t *testing.T) {
	tests := []struct {
		name        string
		startMsg    string
		closure     func() (string, error)
		wantResult  string
		wantErr     bool
	}{
		{
			name:     "successful closure",
			startMsg: "Processing...",
			closure: func() (string, error) {
				return "success", nil
			},
			wantResult: "success",
			wantErr:    false,
		},
		{
			name:     "closure with error",
			startMsg: "Failing...",
			closure: func() (string, error) {
				return "", errors.New("test error")
			},
			wantResult: "",
			wantErr:    true,
		},
		{
			name:     "closure with result and no error",
			startMsg: "Computing...",
			closure: func() (string, error) {
				return "computed value", nil
			},
			wantResult: "computed value",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Spinner(tt.startMsg, tt.closure)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Spinner() error = nil, want error")
				}
			} else {
				if err != nil {
					t.Errorf("Spinner() unexpected error = %v", err)
				}
				if got != tt.wantResult {
					t.Errorf("Spinner() = %q, want %q", got, tt.wantResult)
				}
			}
		})
	}
}

func TestSpinner_GenericTypes(t *testing.T) {
	// Test Spinner with int type
	t.Run("Spinner with int", func(t *testing.T) {
		result, err := Spinner("Counting...", func() (int, error) {
			return 42, nil
		})

		if err != nil {
			t.Errorf("Spinner() unexpected error = %v", err)
		}
		if result != 42 {
			t.Errorf("Spinner() = %d, want 42", result)
		}
	})

	// Test Spinner with slice type
	t.Run("Spinner with slice", func(t *testing.T) {
		result, err := Spinner("Loading...", func() ([]string, error) {
			return []string{"a", "b", "c"}, nil
		})

		if err != nil {
			t.Errorf("Spinner() unexpected error = %v", err)
		}
		if len(result) != 3 {
			t.Errorf("Spinner() returned slice with length %d, want 3", len(result))
		}
	})

	// Test Spinner with struct type
	t.Run("Spinner with struct", func(t *testing.T) {
		type TestStruct struct {
			Name  string
			Value int
		}

		result, err := Spinner("Creating...", func() (TestStruct, error) {
			return TestStruct{Name: "test", Value: 100}, nil
		})

		if err != nil {
			t.Errorf("Spinner() unexpected error = %v", err)
		}
		if result.Name != "test" || result.Value != 100 {
			t.Errorf("Spinner() = %+v, want {Name:test Value:100}", result)
		}
	})
}

func TestAddGroup(t *testing.T) {
	t.Run("adds group with single command", func(t *testing.T) {
		parent := &cobra.Command{Use: "parent"}
		child := &cobra.Command{Use: "child"}

		AddGroup(parent, "Test Group", child)

		// Verify group was added
		foundGroup := false
		for _, g := range parent.Groups() {
			if g.Title == "Test Group" && g.ID == "Test Group" {
				foundGroup = true
				break
			}
		}
		if !foundGroup {
			t.Error("AddGroup() did not add group to parent command")
		}

		// Verify child command was added with correct GroupID
		if child.GroupID != "Test Group" {
			t.Errorf("child.GroupID = %q, want %q", child.GroupID, "Test Group")
		}

		// Verify child was added to parent
		foundChild := false
		for _, cmd := range parent.Commands() {
			if cmd.Use == "child" {
				foundChild = true
				break
			}
		}
		if !foundChild {
			t.Error("AddGroup() did not add child command to parent")
		}
	})

	t.Run("adds group with multiple commands", func(t *testing.T) {
		parent := &cobra.Command{Use: "parent"}
		child1 := &cobra.Command{Use: "child1"}
		child2 := &cobra.Command{Use: "child2"}
		child3 := &cobra.Command{Use: "child3"}

		AddGroup(parent, "Multi Command Group", child1, child2, child3)

		// Verify all children have the same GroupID
		for i, child := range []*cobra.Command{child1, child2, child3} {
			if child.GroupID != "Multi Command Group" {
				t.Errorf("child%d.GroupID = %q, want %q", i+1, child.GroupID, "Multi Command Group")
			}
		}

		// Verify all children were added
		if len(parent.Commands()) != 3 {
			t.Errorf("parent has %d commands, want 3", len(parent.Commands()))
		}
	})
}
