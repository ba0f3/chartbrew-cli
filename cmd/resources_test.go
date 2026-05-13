package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestResourceCommandsRegistered(t *testing.T) {
	resources := []string{"teams", "dashboards", "connections", "datasets", "data-requests", "charts"}
	for _, resource := range resources {
		t.Run(resource, func(t *testing.T) {
			var stdout bytes.Buffer
			root := newRootCommand(&fakeRequester{}, strings.NewReader(""), &stdout, &bytes.Buffer{})
			root.SetArgs([]string{resource, "--help"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			help := stdout.String()
			for _, verb := range []string{"list", "get", "create", "update"} {
				if !strings.Contains(help, verb) {
					t.Fatalf("%s help missing %s:\n%s", resource, verb, help)
				}
			}
		})
	}
}

func TestDocumentedDeleteCommandsRegistered(t *testing.T) {
	resources := []string{"dashboards", "datasets", "charts"}
	for _, resource := range resources {
		t.Run(resource, func(t *testing.T) {
			var stdout bytes.Buffer
			root := newRootCommand(&fakeRequester{}, strings.NewReader(""), &stdout, &bytes.Buffer{})
			root.SetArgs([]string{resource, "--help"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(stdout.String(), "delete") {
				t.Fatalf("%s help missing delete:\n%s", resource, stdout.String())
			}
		})
	}
}

func TestUndocumentedDeleteCommandsNotRegistered(t *testing.T) {
	resources := []string{"teams", "connections", "data-requests"}
	for _, resource := range resources {
		t.Run(resource, func(t *testing.T) {
			var stdout bytes.Buffer
			root := newRootCommand(&fakeRequester{}, strings.NewReader(""), &stdout, &bytes.Buffer{})
			root.SetArgs([]string{resource, "--help"})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			if strings.Contains(stdout.String(), "delete") {
				t.Fatalf("%s help should not include delete:\n%s", resource, stdout.String())
			}
		})
	}
}
