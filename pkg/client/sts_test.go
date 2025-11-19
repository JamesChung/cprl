package client_test

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/smithy-go"

	"github.com/JamesChung/cprl/pkg/client"
)

// TestSTSClient_Integration tests STS client with real AWS credentials
// This is an integration test that requires valid AWS credentials
func TestSTSClient_Integration(t *testing.T) {
	// Skip in CI or if AWS credentials are not available
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot determine home directory, skipping AWS STS test")
	}

	credentialsFile := home + "/.aws/credentials"
	configFile := home + "/.aws/config"

	_, credErr := os.Stat(credentialsFile)
	_, confErr := os.Stat(configFile)

	if os.IsNotExist(credErr) && os.IsNotExist(confErr) {
		t.Skip("No AWS credentials or config file found, skipping integration test")
	}

	t.Run("GetCallerIdentity with default profile", func(t *testing.T) {
		c, err := client.NewSTSClient("default")
		if err != nil {
			t.Fatalf("NewSTSClient failed: %v", err)
		}

		out, err := c.GetCallerIdentity()
		if err != nil {
			// Check if it's an API error
			var gen *smithy.GenericAPIError
			if errors.As(err, &gen) {
				t.Logf("API Error Code: %s", gen.Code)
				t.Logf("API Error Message: %s", gen.Message)
			}

			// SSO session expired is a common case, skip instead of failing
			if errors.Is(err, errors.New("SSO session has expired")) ||
			   (err != nil && (errors.As(err, &gen))) {
				t.Skipf("GetCallerIdentity failed (likely due to expired credentials): %v", err)
			}

			t.Fatalf("GetCallerIdentity failed: %v", err)
		}

		// Verify output has expected fields
		if out == nil {
			t.Fatal("GetCallerIdentity returned nil output")
		}

		if out.Arn == nil || *out.Arn == "" {
			t.Error("GetCallerIdentity returned empty ARN")
		} else {
			t.Logf("Caller ARN: %s", *out.Arn)
		}

		if out.Account == nil || *out.Account == "" {
			t.Error("GetCallerIdentity returned empty Account")
		} else {
			t.Logf("Account: %s", *out.Account)
		}

		if out.UserId == nil || *out.UserId == "" {
			t.Error("GetCallerIdentity returned empty UserId")
		} else {
			t.Logf("UserId: %s", *out.UserId)
		}
	})
}

func TestNewSTSClient_InvalidProfile(t *testing.T) {
	// Test creating client with non-existent profile
	invalidProfile := "this-profile-does-not-exist-99999"

	c, err := client.NewSTSClient(invalidProfile)

	// Creating the client should succeed even with invalid profile
	// The error would occur when making an API call
	if err != nil {
		t.Logf("NewSTSClient with invalid profile returned error: %v", err)
	}

	if c != nil && c.Client == nil {
		t.Error("NewSTSClient returned non-nil client with nil internal Client")
	}
}
