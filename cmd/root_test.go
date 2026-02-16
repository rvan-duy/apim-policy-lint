package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateNoArgsReturnsError(t *testing.T) {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"validate"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no file argument is provided")
	}
}

func TestValidateNonExistentFileReturnsError(t *testing.T) {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"validate", "does-not-exist.xml"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestValidateDirectoryReturnsError(t *testing.T) {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"validate", t.TempDir()})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for directory path")
	}
}

func TestValidateExistingFilePrintsPlaceholder(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.xml")
	if err := os.WriteFile(filePath, []byte("<policies/>"), 0o600); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	rootCmd := NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"validate", filePath})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected success for valid file path, got error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "validate: not implemented yet for "+filePath) {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestRootHelpIncludesValidateCommand(t *testing.T) {
	rootCmd := NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected help command to succeed, got error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "validate") {
		t.Fatalf("expected help output to include validate command, got: %q", output)
	}
}
