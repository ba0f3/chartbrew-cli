package cmd

import (
	"fmt"
	"encoding/json"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		if outputFmt == "json" {
			out, _ := json.MarshalIndent(map[string]string{"version": Version}, "", "  ")
			fmt.Println(string(out))
		} else {
			fmt.Printf("Version: %s\n", Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
