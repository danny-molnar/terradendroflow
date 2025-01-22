package main

import (
	"os"
	"testing"
)

func TestParseAndPrettifyStdout(t *testing.T) {
	inputFile := "testdata/sample_plan.stdout"
	outputFile := "testdata/output.md"

	err := parseAndPrettifyStdout(inputFile, outputFile)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = os.Stat(outputFile)
	if os.IsNotExist(err) {
		t.Fatalf("Expected output file to be created, but it was not")
	}

	// Clean up
	os.Remove(outputFile)
}
