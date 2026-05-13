package cmd

import "testing"

func TestRootCommandMetadata(t *testing.T) {
	root := NewRootCommand()

	if root.Use != "chartbrew" {
		t.Fatalf("Use = %q, want chartbrew", root.Use)
	}
	if root.Short == "" {
		t.Fatal("Short help must not be empty")
	}

	flagNames := []string{"base-url", "token", "config", "output"}
	for _, name := range flagNames {
		if root.PersistentFlags().Lookup(name) == nil {
			t.Fatalf("missing persistent flag %q", name)
		}
	}
}
