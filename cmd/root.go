package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	outputFmt string
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "An Agent-First CLI template",
	Long: `A template for creating CLI tools that are perfectly suited 
for orchestration by AI agents, supporting structured outputs and secure inputs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config, loggers, etc.
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "json", "Output format (json, markdown, raw)")
}
