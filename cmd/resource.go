package cmd

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ba0f3/chartbrew-cli/internal/body"
	"github.com/spf13/cobra"
)

type Route struct {
	Use         string
	Short       string
	Method      string
	Path        func(values map[string]string) string
	IDFlags     []string
	NeedsBody   bool
	Destructive bool
}

func newResourceCommand(name, short string, routes []Route, state *appState) *cobra.Command {
	resource := &cobra.Command{
		Use:   name,
		Short: short,
	}
	for _, route := range routes {
		resource.AddCommand(newRouteCommand(route, state))
	}
	return resource
}

func newRouteCommand(route Route, state *appState) *cobra.Command {
	var data string
	var dataFile string
	cmd := &cobra.Command{
		Use:   route.Use,
		Short: route.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			values := map[string]string{}
			for _, name := range route.IDFlags {
				value, err := requiredFlag(cmd, name)
				if err != nil {
					return err
				}
				values[name] = value
			}

			var reqBody []byte
			if route.NeedsBody {
				var err error
				reqBody, err = body.ReadJSON(body.Source{Data: data, DataFile: dataFile, Stdin: state.stdin})
				if err != nil {
					return err
				}
			}
			if route.Destructive && !state.allowDelete {
				return errors.New("delete commands require allow_delete: true in the Chartbrew config file")
			}

			resp, err := state.api.Do(cmd.Context(), route.Method, route.Path(values), reqBody)
			if err != nil {
				return err
			}
			return writeValue(state, resp)
		},
	}
	for _, name := range route.IDFlags {
		cmd.Flags().String(name, "", "Required "+strings.ReplaceAll(name, "-", " "))
	}
	if route.NeedsBody {
		cmd.Flags().StringVar(&data, "data", "", "Inline JSON request body")
		cmd.Flags().StringVar(&dataFile, "data-file", "", "Path to JSON request body file, or - for stdin")
	}
	return cmd
}

func collectionRoutes(listPath, createPath func(map[string]string) string, flags []string) []Route {
	return []Route{
		{Use: "list", Short: "List resources", Method: http.MethodGet, Path: listPath, IDFlags: flags},
		{Use: "create", Short: "Create a resource", Method: http.MethodPost, Path: createPath, IDFlags: flags, NeedsBody: true},
	}
}
