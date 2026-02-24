package cmd

import (
	"bytes"
	"encoding/json"
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

func TestValidateExistingFilePrintsJSON(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.xml")
	if err := os.WriteFile(filePath, []byte("<policies><inbound><base/></inbound></policies>"), 0o600); err != nil {
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

	var parsed map[string]any
	if err := json.Unmarshal(out.Bytes(), &parsed); err != nil {
		t.Fatalf("expected valid json output, got error: %v. output=%q", err, out.String())
	}

	if parsed["schemaVersion"] != "v1" {
		t.Fatalf("expected schemaVersion=v1, got: %#v", parsed["schemaVersion"])
	}
	if parsed["sourceFile"] != filePath {
		t.Fatalf("expected sourceFile=%q, got: %#v", filePath, parsed["sourceFile"])
	}
	if _, ok := parsed["raw"]; !ok {
		t.Fatalf("expected raw field in output json, got: %#v", parsed)
	}
	if _, ok := parsed["index"]; ok {
		t.Fatalf("did not expect index field in output json, got: %#v", parsed)
	}
}

func TestValidateOutFlagWritesJSONFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.xml")
	outPath := filepath.Join(tempDir, "parsed.json")
	if err := os.WriteFile(filePath, []byte("<policies><inbound><base/></inbound></policies>"), 0o600); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	rootCmd := NewRootCmd()
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stdout)
	rootCmd.SetArgs([]string{"validate", filePath, "--out", outPath})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected success for valid file path with --out, got error: %v", err)
	}

	written, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("expected output file to be created, got error: %v", err)
	}

	if len(written) == 0 {
		t.Fatal("expected output file to contain json")
	}
	if len(stdout.Bytes()) == 0 {
		t.Fatal("expected stdout to still contain json output")
	}
}

func TestValidateMalformedXMLReturnsError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "invalid.xml")
	if err := os.WriteFile(filePath, []byte("<policies><inbound></policies>"), 0o600); err != nil {
		t.Fatalf("failed to create malformed temp file: %v", err)
	}

	rootCmd := NewRootCmd()
	rootCmd.SetArgs([]string{"validate", filePath})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected malformed xml to return error")
	}
	if !strings.Contains(err.Error(), "invalid xml policy file") {
		t.Fatalf("expected parse context in error, got: %v", err)
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
