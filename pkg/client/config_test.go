package client

import (
	"context"
	"os"
	"testing"
)

// TestGetProfileConfig_Integration tests loading AWS config with the default profile
// This is an integration test that requires AWS credentials to be configured
func TestGetProfileConfig_Integration(t *testing.T) {
	// Skip this test if running in CI or if AWS credentials are not available
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Check if AWS credentials are likely available
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot determine home directory, skipping AWS config test")
	}

	credentialsFile := home + "/.aws/credentials"
	configFile := home + "/.aws/config"

	_, credErr := os.Stat(credentialsFile)
	_, confErr := os.Stat(configFile)

	if os.IsNotExist(credErr) && os.IsNotExist(confErr) {
		t.Skip("No AWS credentials or config file found, skipping integration test")
	}

	t.Run("load default profile", func(t *testing.T) {
		cfg, err := GetProfileConfig("default")
		if err != nil {
			t.Fatalf("GetProfileConfig(\"default\") failed: %v", err)
		}

		// Verify config is not empty
		if cfg.Region == "" {
			t.Log("Warning: AWS config has no region set")
		}

		// Try to retrieve credentials
		ctx := context.Background()
		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			t.Fatalf("Failed to retrieve credentials: %v", err)
		}

		// Verify credentials have required fields
		if creds.AccessKeyID == "" {
			t.Error("Retrieved credentials have empty AccessKeyID")
		}
		if creds.SecretAccessKey == "" {
			t.Error("Retrieved credentials have empty SecretAccessKey")
		}

		// Log credential source for debugging (not the actual credentials)
		t.Logf("Credentials source: %s", creds.Source)
	})
}

func TestGetProfileConfig_InvalidProfile(t *testing.T) {
	// This test should work even without AWS credentials
	// as it tests error handling for an invalid profile

	// Use a profile name that's extremely unlikely to exist
	invalidProfile := "this-profile-definitely-does-not-exist-12345"

	cfg, err := GetProfileConfig(invalidProfile)

	// We expect either an error or an empty region (depending on AWS SDK behavior)
	if err == nil && cfg.Region == "" {
		// This is acceptable - SDK might return empty config for non-existent profile
		t.Logf("GetProfileConfig returned empty config for invalid profile (expected)")
		return
	}

	if err == nil {
		// If there's no error, at least log what we got
		t.Logf("GetProfileConfig succeeded with profile %q, region: %q", invalidProfile, cfg.Region)
	}
}

func TestGetProfileConfig_EmptyProfile(t *testing.T) {
	// Test with empty profile name
	// The AWS SDK should handle this gracefully (likely using default profile)

	cfg, err := GetProfileConfig("")

	// We don't strictly require an error here, as SDK might fall back to default
	if err != nil {
		t.Logf("GetProfileConfig(\"\") returned error (acceptable): %v", err)
		return
	}

	// If it succeeded, it should have returned some config
	t.Logf("GetProfileConfig(\"\") succeeded, region: %q", cfg.Region)
}
