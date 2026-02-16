package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "apimlint",
		Short: "Lint APIM policies",
	}

	rootCmd.AddCommand(newValidateCmd())

	return rootCmd
}

func Execute() error {
	return NewRootCmd().Execute()
}
