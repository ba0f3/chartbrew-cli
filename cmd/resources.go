package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
)

func addResourceCommands(root *cobra.Command, state *appState) {
	root.AddCommand(
		newTeamsCommand(state),
		newDashboardsCommand(state),
		newConnectionsCommand(state),
		newDatasetsCommand(state),
		newDataRequestsCommand(state),
		newChartsCommand(state),
	)
}

func newTeamsCommand(state *appState) *cobra.Command {
	return newResourceCommand("teams", "Manage Chartbrew teams", []Route{
		{Use: "list", Short: "List teams", Method: http.MethodGet, Path: func(map[string]string) string { return "/team" }},
		{Use: "get", Short: "Get a team", Method: http.MethodGet, Path: func(v map[string]string) string { return "/team/" + v["team-id"] }, IDFlags: []string{"team-id"}},
		{Use: "create", Short: "Create a team", Method: http.MethodPost, Path: func(map[string]string) string { return "/team" }, NeedsBody: true},
		{Use: "update", Short: "Update a team", Method: http.MethodPut, Path: func(v map[string]string) string { return "/team/" + v["team-id"] }, IDFlags: []string{"team-id"}, NeedsBody: true},
	}, state)
}

func newDashboardsCommand(state *appState) *cobra.Command {
	return newResourceCommand("dashboards", "Manage Chartbrew dashboards", []Route{
		{Use: "list", Short: "List dashboards", Method: http.MethodGet, Path: func(v map[string]string) string { return "/project/team/" + v["team-id"] }, IDFlags: []string{"team-id"}},
		{Use: "get", Short: "Get a dashboard", Method: http.MethodGet, Path: func(v map[string]string) string { return "/project/" + v["dashboard-id"] }, IDFlags: []string{"dashboard-id"}},
		{Use: "create", Short: "Create a dashboard", Method: http.MethodPost, Path: func(map[string]string) string { return "/project" }, NeedsBody: true},
		{Use: "update", Short: "Update a dashboard", Method: http.MethodPut, Path: func(v map[string]string) string { return "/project/" + v["dashboard-id"] }, IDFlags: []string{"dashboard-id"}, NeedsBody: true},
	}, state)
}

func newConnectionsCommand(state *appState) *cobra.Command {
	return newResourceCommand("connections", "Manage Chartbrew connections", []Route{
		{Use: "list", Short: "List connections", Method: http.MethodGet, Path: func(v map[string]string) string { return "/team/" + v["team-id"] + "/connections" }, IDFlags: []string{"team-id"}},
		{Use: "get", Short: "Get a connection", Method: http.MethodGet, Path: func(v map[string]string) string {
			return "/team/" + v["team-id"] + "/connections/" + v["connection-id"]
		}, IDFlags: []string{"team-id", "connection-id"}},
		{Use: "create", Short: "Create a connection", Method: http.MethodPost, Path: func(v map[string]string) string { return "/team/" + v["team-id"] + "/connections" }, IDFlags: []string{"team-id"}, NeedsBody: true},
		{Use: "update", Short: "Update a connection", Method: http.MethodPut, Path: func(v map[string]string) string {
			return "/team/" + v["team-id"] + "/connections/" + v["connection-id"]
		}, IDFlags: []string{"team-id", "connection-id"}, NeedsBody: true},
	}, state)
}

func newDatasetsCommand(state *appState) *cobra.Command {
	return newResourceCommand("datasets", "Manage Chartbrew datasets", []Route{
		{Use: "list", Short: "List datasets", Method: http.MethodGet, Path: func(v map[string]string) string { return "/team/" + v["team-id"] + "/datasets" }, IDFlags: []string{"team-id"}},
		{Use: "get", Short: "Get a dataset", Method: http.MethodGet, Path: func(v map[string]string) string { return "/team/" + v["team-id"] + "/datasets/" + v["dataset-id"] }, IDFlags: []string{"team-id", "dataset-id"}},
		{Use: "create", Short: "Create a dataset", Method: http.MethodPost, Path: func(v map[string]string) string { return "/team/" + v["team-id"] + "/datasets" }, IDFlags: []string{"team-id"}, NeedsBody: true},
		{Use: "update", Short: "Update a dataset", Method: http.MethodPut, Path: func(v map[string]string) string { return "/team/" + v["team-id"] + "/datasets/" + v["dataset-id"] }, IDFlags: []string{"team-id", "dataset-id"}, NeedsBody: true},
	}, state)
}

func newDataRequestsCommand(state *appState) *cobra.Command {
	return newResourceCommand("data-requests", "Manage Chartbrew data requests", []Route{
		{Use: "list", Short: "List data requests", Method: http.MethodGet, Path: func(v map[string]string) string {
			return "/team/" + v["team-id"] + "/datasets/" + v["dataset-id"] + "/dataRequests"
		}, IDFlags: []string{"team-id", "dataset-id"}},
		{Use: "get", Short: "Get a data request", Method: http.MethodGet, Path: func(v map[string]string) string {
			return "/team/" + v["team-id"] + "/datasets/" + v["dataset-id"] + "/dataRequests/" + v["request-id"]
		}, IDFlags: []string{"team-id", "dataset-id", "request-id"}},
		{Use: "create", Short: "Create a data request", Method: http.MethodPost, Path: func(v map[string]string) string {
			return "/team/" + v["team-id"] + "/datasets/" + v["dataset-id"] + "/dataRequests"
		}, IDFlags: []string{"team-id", "dataset-id"}, NeedsBody: true},
		{Use: "update", Short: "Update a data request", Method: http.MethodPut, Path: func(v map[string]string) string {
			return "/team/" + v["team-id"] + "/datasets/" + v["dataset-id"] + "/dataRequests/" + v["request-id"]
		}, IDFlags: []string{"team-id", "dataset-id", "request-id"}, NeedsBody: true},
	}, state)
}

func newChartsCommand(state *appState) *cobra.Command {
	return newResourceCommand("charts", "Manage Chartbrew charts", []Route{
		{Use: "list", Short: "List charts", Method: http.MethodGet, Path: func(v map[string]string) string { return "/project/" + v["dashboard-id"] + "/chart" }, IDFlags: []string{"dashboard-id"}},
		{Use: "get", Short: "Get a chart", Method: http.MethodGet, Path: func(v map[string]string) string { return "/project/" + v["dashboard-id"] + "/chart/" + v["chart-id"] }, IDFlags: []string{"dashboard-id", "chart-id"}},
		{Use: "create", Short: "Create a chart", Method: http.MethodPost, Path: func(v map[string]string) string { return "/project/" + v["dashboard-id"] + "/chart" }, IDFlags: []string{"dashboard-id"}, NeedsBody: true},
		{Use: "update", Short: "Update a chart", Method: http.MethodPut, Path: func(v map[string]string) string { return "/project/" + v["dashboard-id"] + "/chart/" + v["chart-id"] }, IDFlags: []string{"dashboard-id", "chart-id"}, NeedsBody: true},
	}, state)
}
