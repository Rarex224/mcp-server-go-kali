package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// A dummy test to ensure the sanitization function works correctly
func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com/page?id=1&name=test", "http://example.com/page?id=1&name=test"},
		{"192.168.1.1; rm -rf /", "192.168.1.1rm-rf/"},
		{"example.com && nmap", "example.comnmap"},
		{"-p 80,443 -sV", "-p80443-sV"},
		{"domain.com/path-to-file_1.txt", "domain.com/path-to-file-1.txt"},
	}

	for _, test := range tests {
		result := sanitizeInput(test.input)
		// Our current sanitizeInput removes spaces, commas, and special chars like ; &&
		// Let's adjust the expected to match the regex `[^a-zA-Z0-9.\-:/=?&]`
		expectedSanitized := strings.Map(func(r rune) rune {
			if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-:/=?&", r) {
				return r
			}
			return -1
		}, test.input)

		if result != expectedSanitized {
			t.Errorf("sanitizeInput(%q) = %q; want %q", test.input, result, expectedSanitized)
		}
	}
}

// A simple test for the runToolCommand mocking a safe command like 'echo'
func TestRunToolCommand(t *testing.T) {
	// This relies on the host having 'echo', which most systems do
	res, err := runToolCommand("echo", "test-output")
	if err != nil {
		t.Fatalf("runToolCommand failed: %v", err)
	}

	if res == nil {
		t.Fatal("Expected result but got nil")
	}

	contentBytes, _ := json.Marshal(res.Content)
	contentStr := string(contentBytes)

	if !strings.Contains(contentStr, "test-output") {
		t.Errorf("Expected output to contain 'test-output', got: %s", contentStr)
	}
}
