package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

type fakeRequester struct {
	method string
	path   string
	body   string
	err    error
}

func (f *fakeRequester) Do(ctx context.Context, method, path string, body []byte) (json.RawMessage, error) {
	f.method = method
	f.path = path
	f.body = string(body)
	if f.err != nil {
		return nil, f.err
	}
	return json.RawMessage(`{"ok":true}`), nil
}

func executeTestCommand(api *fakeRequester, stdin string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := newRootCommand(api, strings.NewReader(stdin), &stdout, &stderr)
	root.SetArgs(args)
	err := root.Execute()
	return stdout.String(), stderr.String(), err
}

func withConfig(args ...string) []string {
	base := []string{"--base-url", "https://api.example", "--token", "token"}
	return append(base, args...)
}

func TestListCommandCallsExpectedPath(t *testing.T) {
	api := &fakeRequester{}
	_, _, err := executeTestCommand(api, "", withConfig("datasets", "list", "--team-id", "7")...)
	if err != nil {
		t.Fatal(err)
	}
	if api.method != http.MethodGet || api.path != "/team/7/datasets" {
		t.Fatalf("request = %s %s", api.method, api.path)
	}
}

func TestGetCommandRequiresID(t *testing.T) {
	api := &fakeRequester{}
	_, _, err := executeTestCommand(api, "", withConfig("datasets", "get", "--team-id", "7")...)
	if err == nil {
		t.Fatal("expected missing dataset-id error")
	}
	if !strings.Contains(err.Error(), "--dataset-id") {
		t.Fatalf("error = %v", err)
	}
}

func TestCreateCommandReadsData(t *testing.T) {
	api := &fakeRequester{}
	_, _, err := executeTestCommand(api, "", withConfig("teams", "create", "--data", `{"name":"demo"}`)...)
	if err != nil {
		t.Fatal(err)
	}
	if api.method != http.MethodPost || api.path != "/team" {
		t.Fatalf("request = %s %s", api.method, api.path)
	}
	if api.body != `{"name":"demo"}` {
		t.Fatalf("body = %s", api.body)
	}
}

func TestUpdateCommandReadsDataAndID(t *testing.T) {
	api := &fakeRequester{}
	_, _, err := executeTestCommand(api, `{"name":"demo"}`, withConfig("teams", "update", "--team-id", "9", "--data-file", "-")...)
	if err != nil {
		t.Fatal(err)
	}
	if api.method != http.MethodPut || api.path != "/team/9" {
		t.Fatalf("request = %s %s", api.method, api.path)
	}
	if api.body != `{"name":"demo"}` {
		t.Fatalf("body = %s", api.body)
	}
}

func TestMissingConfigReturnsError(t *testing.T) {
	api := &fakeRequester{}
	_, _, err := executeTestCommand(api, "", "teams", "list")
	if err == nil {
		t.Fatal("expected missing config error")
	}
}

func TestInvalidOutputFormatFailsBeforeRequest(t *testing.T) {
	api := &fakeRequester{}
	_, _, err := executeTestCommand(api, "", "--base-url", "https://api.example", "--token", "token", "--output", "xml", "teams", "list")
	if err == nil {
		t.Fatal("expected output format error")
	}
	if api.method != "" {
		t.Fatalf("request should not run, got %s", api.method)
	}
}

func TestCommandErrorDoesNotWritePartialSuccessToStdout(t *testing.T) {
	api := &fakeRequester{err: errors.New("boom")}
	stdout, _, err := executeTestCommand(api, "", withConfig("teams", "list")...)
	if err == nil {
		t.Fatal("expected command error")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q", stdout)
	}
}
