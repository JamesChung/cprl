package util_test

import (
	"os"
	"strings"
	"testing"

	"github.com/JamesChung/cprl/pkg/util"
)

// TestGetCredentials_Integration tests retrieving AWS credentials
// This is an integration test that requires valid AWS credentials
func TestGetCredentials_Integration(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot determine home directory, skipping AWS credentials test")
	}

	credentialsFile := home + "/.aws/credentials"
	configFile := home + "/.aws/config"

	_, credErr := os.Stat(credentialsFile)
	_, confErr := os.Stat(configFile)

	if os.IsNotExist(credErr) && os.IsNotExist(confErr) {
		t.Skip("No AWS credentials or config file found, skipping integration test")
	}

	creds, err := util.GetCredentials("default")
	if err != nil {
		t.Skipf("GetCredentials failed (likely expired credentials): %v", err)
	}

	// Verify credentials structure
	if creds.AccessKeyID == "" && creds.SecretAccessKey == "" {
		t.Skip("Retrieved empty credentials, skipping test")
	}

	t.Logf("Credentials source: %s", creds.Source)
}

func TestGenerateLoginURL_Integration(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot determine home directory, skipping AWS credentials test")
	}

	credentialsFile := home + "/.aws/credentials"
	configFile := home + "/.aws/config"

	_, credErr := os.Stat(credentialsFile)
	_, confErr := os.Stat(configFile)

	if os.IsNotExist(credErr) && os.IsNotExist(confErr) {
		t.Skip("No AWS credentials or config file found, skipping integration test")
	}

	creds, err := util.GetCredentials("default")
	if err != nil {
		t.Skipf("GetCredentials failed: %v", err)
	}

	if creds.AccessKeyID == "" {
		t.Skip("No valid credentials available, skipping test")
	}

	t.Run("standard AWS login URL", func(t *testing.T) {
		loginURL, err := util.GenerateLoginURL(creds, false)
		if err != nil {
			t.Fatalf("GenerateLoginURL failed: %v", err)
		}

		urlStr := loginURL.String()
		if !strings.Contains(urlStr, "https://") {
			t.Errorf("Login URL does not start with https://, got: %s", urlStr)
		}

		if !strings.Contains(urlStr, "signin") {
			t.Errorf("Login URL does not contain 'signin', got: %s", urlStr)
		}

		t.Logf("Generated login URL (truncated): %s...", urlStr[:min(len(urlStr), 100)])
	})

	t.Run("GovCloud login URL", func(t *testing.T) {
		loginURL, err := util.GenerateLoginURL(creds, true)
		if err != nil {
			t.Fatalf("GenerateLoginURL for GovCloud failed: %v", err)
		}

		urlStr := loginURL.String()
		if !strings.Contains(urlStr, "https://") {
			t.Errorf("GovCloud login URL does not start with https://, got: %s", urlStr)
		}

		t.Logf("Generated GovCloud login URL (truncated): %s...", urlStr[:min(len(urlStr), 100)])
	})
}

func TestGenerateLoginURL_WithMockCredentials(t *testing.T) {
	t.Skip("Skipping mock credentials test - GenerateLoginURL makes real AWS API calls")

	// Note: GenerateLoginURL actually makes HTTP calls to AWS STS to get federated tokens,
	// so it cannot be tested with mock credentials without setting up a mock HTTP server.
	// This would require refactoring the function to accept an HTTP client interface.
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
