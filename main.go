package main

import (
	"os"

	"github.com/username/agent-cli-template/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// You can customize the error handling and exit codes based on error type
		os.Exit(1)
	}
}
