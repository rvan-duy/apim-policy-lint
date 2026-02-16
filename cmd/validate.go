package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
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

			cmd.Printf("validate: not implemented yet for %s\n", filePath)
			return nil
		},
	}
}
