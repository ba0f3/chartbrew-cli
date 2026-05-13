package cmd

import "github.com/spf13/cobra"

func newVersionCommand(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:         "version",
		Short:       "Print the version number",
		Annotations: map[string]string{"skip_config": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeValue(state, map[string]string{"version": Version})
		},
	}
}
