package main

import (
	"os"

	"github.com/ba0f3/chartbrew-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(cmd.HandleError(err))
	}
}
