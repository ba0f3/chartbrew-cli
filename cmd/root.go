package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ba0f3/chartbrew-cli/internal/client"
	"github.com/ba0f3/chartbrew-cli/internal/config"
	"github.com/ba0f3/chartbrew-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

type requester interface {
	Do(ctx context.Context, method, path string, body []byte) (json.RawMessage, error)
}

type appState struct {
	baseURL    string
	token      string
	configPath string
	outputFmt  string
	api        requester
	stdin      io.Reader
	stdout     io.Writer
	stderr     io.Writer
}

func Execute() error {
	return NewRootCommand().Execute()
}

func NewRootCommand() *cobra.Command {
	return newRootCommand(nil, os.Stdin, os.Stdout, os.Stderr)
}

func newRootCommand(api requester, stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	state := &appState{
		outputFmt:  "json",
		configPath: config.DefaultConfigPath(),
		api:        api,
		stdin:      stdin,
		stdout:     stdout,
		stderr:     stderr,
	}

	root := &cobra.Command{
		Use:   "chartbrew",
		Short: "Agent-first CLI for the Chartbrew API",
		Long:  "Manage Chartbrew teams, dashboards, connections, datasets, data requests, and charts from scripts and AI agents.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if skipConfig(cmd) {
				return nil
			}
			cfg, err := config.Resolve(config.Sources{
				FlagBaseURL: state.baseURL,
				FlagToken:   state.token,
				FlagOutput:  state.outputFmt,
				DotenvPath:  ".env",
				ConfigPath:  state.configPath,
			})
			if err != nil {
				return err
			}
			if err := cfg.Validate(); err != nil {
				return err
			}
			state.outputFmt = cfg.Output
			if state.api == nil {
				state.api = client.New(cfg.BaseURL, cfg.Token, http.DefaultClient)
			}
			return nil
		},
	}
	root.SilenceUsage = true
	root.SilenceErrors = true
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetIn(stdin)

	root.PersistentFlags().StringVar(&state.baseURL, "base-url", "", "Chartbrew API base URL")
	root.PersistentFlags().StringVar(&state.token, "token", "", "Chartbrew API token; prefer CHARTBREW_TOKEN or config file for secrets")
	root.PersistentFlags().StringVar(&state.configPath, "config", state.configPath, "Path to Chartbrew config JSON")
	root.PersistentFlags().StringVarP(&state.outputFmt, "output", "o", "json", "Output format (json, markdown, raw)")

	root.AddCommand(newVersionCommand(state))
	addResourceCommands(root, state)
	return root
}

func HandleError(err error) int {
	if err == nil {
		return 0
	}
	code := 1
	var httpErr client.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode > 0 {
		code = httpErr.StatusCode
	}
	_ = output.WriteError(os.Stderr, code, err.Error())
	return code
}

func skipConfig(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		if current.Annotations["skip_config"] == "true" {
			return true
		}
	}
	return false
}

func writeValue(state *appState, value any) error {
	return output.Write(state.stdout, output.Format(state.outputFmt), value)
}

func requiredFlag(cmd *cobra.Command, name string) (string, error) {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("missing required flag --%s", name)
	}
	return value, nil
}
