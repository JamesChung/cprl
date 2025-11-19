package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiff_IdenticalFiles(t *testing.T) {
	old := []byte("line1\nline2\nline3\n")
	new := []byte("line1\nline2\nline3\n")

	result := Diff("old.txt", old, "new.txt", new)

	if result != nil {
		t.Errorf("Diff() of identical files = %q, want nil", result)
	}
}

func TestDiff_SimpleAddition(t *testing.T) {
	old := []byte("line1\nline2\n")
	new := []byte("line1\nline2\nline3\n")

	result := Diff("old.txt", old, "new.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// Check for diff header
	if !strings.Contains(output, "diff old.txt new.txt") {
		t.Errorf("output missing diff header, got:\n%s", output)
	}
	if !strings.Contains(output, "--- old.txt") {
		t.Errorf("output missing old file marker, got:\n%s", output)
	}
	if !strings.Contains(output, "+++ new.txt") {
		t.Errorf("output missing new file marker, got:\n%s", output)
	}

	// Check for added line
	if !strings.Contains(output, "+line3") {
		t.Errorf("output missing added line, got:\n%s", output)
	}
}

func TestDiff_SimpleDeletion(t *testing.T) {
	old := []byte("line1\nline2\nline3\n")
	new := []byte("line1\nline3\n")

	result := Diff("old.txt", old, "new.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// Check for deleted line
	if !strings.Contains(output, "-line2") {
		t.Errorf("output missing deleted line, got:\n%s", output)
	}
}

func TestDiff_Modification(t *testing.T) {
	old := []byte("line1\nline2\nline3\n")
	new := []byte("line1\nmodified line2\nline3\n")

	result := Diff("old.txt", old, "new.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// Modification shows as deletion + addition
	if !strings.Contains(output, "-line2") {
		t.Errorf("output missing original line, got:\n%s", output)
	}
	if !strings.Contains(output, "+modified line2") {
		t.Errorf("output missing modified line, got:\n%s", output)
	}
}

func TestDiff_EmptyToContent(t *testing.T) {
	old := []byte("")
	new := []byte("new line\n")

	result := Diff("empty.txt", old, "content.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	if !strings.Contains(output, "+new line") {
		t.Errorf("output missing added line, got:\n%s", output)
	}
}

func TestDiff_ContentToEmpty(t *testing.T) {
	old := []byte("old line\n")
	new := []byte("")

	result := Diff("content.txt", old, "empty.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	if !strings.Contains(output, "-old line") {
		t.Errorf("output missing deleted line, got:\n%s", output)
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	old := []byte("line1\nline2\nline3\nline4\nline5\n")
	new := []byte("line1\nmodified2\nline3\nline5\nline6\n")

	result := Diff("old.txt", old, "new.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// line2 modified
	if !strings.Contains(output, "-line2") {
		t.Errorf("output missing original line2")
	}
	if !strings.Contains(output, "+modified2") {
		t.Errorf("output missing modified line2")
	}

	// line4 deleted
	if !strings.Contains(output, "-line4") {
		t.Errorf("output missing deleted line4")
	}

	// line6 added
	if !strings.Contains(output, "+line6") {
		t.Errorf("output missing added line6")
	}
}

func TestDiff_NoTrailingNewline(t *testing.T) {
	old := []byte("line1\nline2")
	new := []byte("line1\nline2\n")

	result := Diff("old.txt", old, "new.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// Should mention missing newline
	if !strings.Contains(output, "No newline at end of file") {
		t.Errorf("output missing newline warning, got:\n%s", output)
	}
}

func TestDiff_BinaryContent(t *testing.T) {
	// Test with binary-like content (null bytes)
	old := []byte("text\x00binary\nmore text\n")
	new := []byte("text\x00different\nmore text\n")

	result := Diff("old.bin", old, "new.bin", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	// Diff should still work with binary content
	output := string(result)
	if !strings.Contains(output, "diff old.bin new.bin") {
		t.Errorf("output missing diff header for binary files")
	}
}

func TestDiff_LargeFile(t *testing.T) {
	// Create a larger file to test performance
	var oldBuf, newBuf bytes.Buffer
	for i := 0; i < 1000; i++ {
		oldBuf.WriteString("line ")
		oldBuf.WriteString(string(rune('0' + (i % 10))))
		oldBuf.WriteString("\n")

		newBuf.WriteString("line ")
		newBuf.WriteString(string(rune('0' + (i % 10))))
		newBuf.WriteString("\n")
	}

	// Modify a few lines
	oldBytes := oldBuf.Bytes()
	newBytes := newBuf.Bytes()
	// Change line 500 (approximately)
	newBytes = bytes.Replace(newBytes, []byte("line 5\n"), []byte("modified 5\n"), 1)

	result := Diff("old.txt", oldBytes, "new.txt", newBytes)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)
	if !strings.Contains(output, "+modified 5") {
		t.Errorf("output missing modification in large file")
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []string
	}{
		{
			name:  "simple lines with trailing newline",
			input: []byte("line1\nline2\nline3\n"),
			want:  []string{"line1\n", "line2\n", "line3\n"},
		},
		{
			name:  "no trailing newline",
			input: []byte("line1\nline2"),
			want:  []string{"line1\n", "line2\n\\ No newline at end of file\n"},
		},
		{
			name:  "single line with newline",
			input: []byte("single\n"),
			want:  []string{"single\n"},
		},
		{
			name:  "single line without newline",
			input: []byte("single"),
			want:  []string{"single\n\\ No newline at end of file\n"},
		},
		{
			name:  "empty file",
			input: []byte(""),
			want:  []string{},
		},
		{
			name:  "just newline",
			input: []byte("\n"),
			want:  []string{"\n"},
		},
		{
			name:  "multiple empty lines",
			input: []byte("\n\n\n"),
			want:  []string{"\n", "\n", "\n"},
		},
		{
			name:  "mixed content",
			input: []byte("start\n\nmiddle\n\nend\n"),
			want:  []string{"start\n", "\n", "middle\n", "\n", "end\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lines(tt.input)

			if len(got) != len(tt.want) {
				t.Errorf("lines() returned %d lines, want %d", len(got), len(tt.want))
				t.Errorf("got: %v", got)
				t.Errorf("want: %v", tt.want)
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("lines()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestDiff_RealWorldExample(t *testing.T) {
	// Simulate a real-world code change
	old := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	x := 42
	fmt.Println(x)
}
`)

	new := []byte(`package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")
	x := 100
	fmt.Println(x)
	os.Exit(0)
}
`)

	result := Diff("old.go", old, "new.go", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// Check for import change
	if !strings.Contains(output, `-import "fmt"`) {
		t.Errorf("output missing old import statement")
	}

	// Check for value change
	if !strings.Contains(output, "-\tx := 42") {
		t.Errorf("output missing old value assignment")
	}
	if !strings.Contains(output, "+\tx := 100") {
		t.Errorf("output missing new value assignment")
	}

	// Check for new line
	if !strings.Contains(output, "+\tos.Exit(0)") {
		t.Errorf("output missing new os.Exit line")
	}
}

func TestDiff_UnifiedFormat(t *testing.T) {
	old := []byte("a\nb\nc\nd\ne\n")
	new := []byte("a\nB\nc\nd\ne\n")

	result := Diff("test.txt", old, "test.txt", new)

	if result == nil {
		t.Fatal("Diff() returned nil, want diff output")
	}

	output := string(result)

	// Check for unified diff format hunk header
	if !strings.Contains(output, "@@") {
		t.Errorf("output missing unified diff hunk header (@@), got:\n%s", output)
	}

	// Should show context lines with space prefix
	lines := strings.Split(output, "\n")
	hasContextLine := false
	for _, line := range lines {
		if strings.HasPrefix(line, " a") || strings.HasPrefix(line, " c") {
			hasContextLine = true
			break
		}
	}
	if !hasContextLine {
		t.Errorf("output missing context lines (space prefix), got:\n%s", output)
	}
}
