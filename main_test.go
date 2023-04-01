package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetIgnoreList(t *testing.T) {
	ignoreFilePath := "test.gptignore"
	ignoreList := []string{"*.txt", "*.log"}

	file, err := os.Create(ignoreFilePath)
	if err != nil {
		t.Fatalf("Failed to create test ignore file: %v", err)
	}
	defer os.Remove(ignoreFilePath)

	for _, pattern := range ignoreList {
		_, err = file.WriteString(pattern + "\n")
		if err != nil {
			t.Fatalf("Failed to write test ignore file: %v", err)
		}
	}
	file.Close()

	result, err := getIgnoreList(ignoreFilePath)
	if err != nil {
		t.Fatalf("Failed to get test ignore file: %v", err)
	}
	if len(result) != len(ignoreList) {
		t.Fatalf("Expected %d patterns, got %d", len(ignoreList), len(result))
	}

	for i, pattern := range ignoreList {
		if result[i] != pattern {
			t.Fatalf("Expected pattern %s, got %s", pattern, result[i])
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	ignoreList := []string{"*.txt", "*.log"}

	testCases := []struct {
		filePath string
		expected bool
	}{
		{"file.txt", true},
		{"file.log", true},
		{"file.go", false},
		{"file.jpg", false},
	}

	for _, tc := range testCases {
		result := shouldIgnore(tc.filePath, ignoreList)
		if result != tc.expected {
			t.Fatalf("Expected %t for file %s, got %t", tc.expected, tc.filePath, result)
		}
	}
}

// Test 1: processDirectory function
func TestProcessDirectory(t *testing.T) {
	// Create a temporary directory and files
	tempDir, err := os.MkdirTemp("", "totxt-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Process the temporary directory with an empty ignore list
	outFile, err := os.CreateTemp("", "totxt-test-output")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(outFile.Name())

	err = processDirectory(tempDir, []string{}, outFile)
	if err != nil {
		t.Fatalf("Failed to process temp directory: %v", err)
	}

	// Check if the output contains the test file content
	output, err := os.ReadFile(outFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp output file: %v", err)
	}
	if !strings.Contains(string(output), "testfile.txt") || !strings.Contains(string(output), "test content") {
		t.Fatalf("Output file does not contain the expected content")
	}
}

// Test 2: preamble file processing
func TestPreambleFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "totxt-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary preamble file
	tempPreambleFile, err := os.CreateTemp("", "totxt-test-preamble")
	if err != nil {
		t.Fatalf("Failed to create temp preamble file: %v", err)
	}
	defer os.Remove(tempPreambleFile.Name())

	preambleContent := "This is a test preamble."
	_, err = tempPreambleFile.WriteString(preambleContent)
	if err != nil {
		t.Fatalf("Failed to write temp preamble file: %v", err)
	}
	tempPreambleFile.Close()

	// Use the temporary preamble file and check if it is added to the output
	tempOutputFile, err := os.CreateTemp("", "totxt-test-output")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutputFile.Name())

	// Write the preamble to the output file
	_, _ = tempOutputFile.WriteString(preambleContent + "\n")

	// Process the temporary directory with an empty ignore list
	err = processDirectory(tempDir, []string{}, tempOutputFile)
	if err != nil {
		t.Fatalf("Failed to process temp directory: %v", err)
	}

	// Check if the output contains the expected preamble content
	output, err := os.ReadFile(tempOutputFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp output file: %v", err)
	}
	if !strings.Contains(string(output), preambleContent) {
		t.Fatalf("Output file does not contain the expected preamble content")
	}
}

// Test 3: error handling
func TestErrorHandling(t *testing.T) {
	// Invalid directory path
	_, err := getIgnoreList("non_existent_directory/.totxtignore")
	if err == nil {
		t.Fatalf("Expected error for non-existent directory, but got nil")
	}

	// Invalid preamble file path
	tempOutputFile, err := os.CreateTemp("", "totxt-test-output")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tempOutputFile.Name())

	err = processDirectory("", []string{}, tempOutputFile)
	if err == nil {
		t.Fatalf("Expected error for non-existent preamble file, but got nil")
	}

	// Invalid output file path
	err = processDirectory("", []string{}, tempOutputFile)
	if err == nil {
		t.Fatalf("Expected error for non-existent output file, but got nil")
	}
}
