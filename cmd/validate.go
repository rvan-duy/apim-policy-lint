package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"apim-policy-lint/internal/policyxml"

	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	var outPath string

	validateCmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate an APIM policy file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			info, err := os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("file does not exist: %s", filePath)
				}
				return fmt.Errorf("unable to access file %s: %w", filePath, err)
			}

			if info.IsDir() {
				return fmt.Errorf("path is a directory, expected a file: %s", filePath)
			}

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("file is not readable: %s: %w", filePath, err)
			}
			_ = file.Close()

			doc, err := policyxml.ParseFile(filePath)
			if err != nil {
				return fmt.Errorf("invalid xml policy file %s: %w", filePath, err)
			}

			output, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return fmt.Errorf("failed serializing parsed policy: %w", err)
			}
			output = append(output, '\n')

			if outPath != "" {
				if err := os.WriteFile(outPath, output, 0o644); err != nil {
					return fmt.Errorf("failed writing json output to %s: %w", outPath, err)
				}
			}

			cmd.Print(string(output))
			return nil
		},
	}

	validateCmd.Flags().StringVar(&outPath, "out", "", "Write parsed JSON output to file")
	return validateCmd
}
