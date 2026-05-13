package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePriorityFlagsEnvDotenvFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	dotenvPath := filepath.Join(dir, ".env")

	if err := os.WriteFile(configPath, []byte(`{"base_url":"https://file.example","token":"file-token","output":"raw"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dotenvPath, []byte("CHARTBREW_API_URL=https://dotenv.example\nCHARTBREW_TOKEN=dotenv-token\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Resolve(Sources{
		FlagBaseURL: "https://flag.example",
		FlagToken:   "flag-token",
		FlagOutput:  "markdown",
		Env: map[string]string{
			EnvBaseURL: "https://env.example",
			EnvToken:   "env-token",
		},
		DotenvPath: dotenvPath,
		ConfigPath: configPath,
	})
	if err != nil {
		t.Fatal(err)
	}

	if cfg.BaseURL != "https://flag.example" {
		t.Fatalf("BaseURL = %q", cfg.BaseURL)
	}
	if cfg.Token != "flag-token" {
		t.Fatalf("Token = %q", cfg.Token)
	}
	if cfg.Output != "markdown" {
		t.Fatalf("Output = %q", cfg.Output)
	}
}

func TestResolveLowerPrioritySources(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	dotenvPath := filepath.Join(dir, ".env")

	if err := os.WriteFile(configPath, []byte(`{"base_url":"https://file.example","token":"file-token"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Resolve(Sources{Env: map[string]string{}, DotenvPath: dotenvPath, ConfigPath: configPath})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://file.example" || cfg.Token != "file-token" {
		t.Fatalf("config file not applied: %+v", cfg)
	}

	if err := os.WriteFile(dotenvPath, []byte("CHARTBREW_API_URL=https://dotenv.example\nCHARTBREW_TOKEN=dotenv-token\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err = Resolve(Sources{Env: map[string]string{}, DotenvPath: dotenvPath, ConfigPath: configPath})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://dotenv.example" || cfg.Token != "dotenv-token" {
		t.Fatalf(".env not applied over config file: %+v", cfg)
	}

	cfg, err = Resolve(Sources{
		Env: map[string]string{
			EnvBaseURL: "https://env.example",
			EnvToken:   "env-token",
		},
		DotenvPath: dotenvPath,
		ConfigPath: configPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://env.example" || cfg.Token != "env-token" {
		t.Fatalf("env not applied over .env: %+v", cfg)
	}
}

func TestResolveAllowDeleteFromConfigFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	if err := os.WriteFile(configPath, []byte(`{"base_url":"https://file.example","token":"file-token","allow_delete":true}`), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Resolve(Sources{Env: map[string]string{}, ConfigPath: configPath})
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.AllowDelete {
		t.Fatal("AllowDelete = false, want true")
	}
}

func TestResolveAllowDeleteCannotBeEnabledByEnvOrFlags(t *testing.T) {
	cfg, err := Resolve(Sources{
		FlagBaseURL: "https://flag.example",
		FlagToken:   "flag-token",
		Env: map[string]string{
			EnvBaseURL:         "https://env.example",
			EnvToken:           "env-token",
			"ALLOW_DELETE":     "true",
			"allow_delete":     "true",
			"CHARTBREW_DELETE": "true",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AllowDelete {
		t.Fatal("AllowDelete = true, want false")
	}
}

func TestValidateRequiresBaseURLAndToken(t *testing.T) {
	if err := (Config{Token: "token"}).Validate(); err == nil {
		t.Fatal("expected missing base URL error")
	}
	if err := (Config{BaseURL: "https://api.chartbrew.com"}).Validate(); err == nil {
		t.Fatal("expected missing token error")
	}
	if err := (Config{BaseURL: "https://api.chartbrew.com", Token: "token", Output: "xml"}).Validate(); err == nil {
		t.Fatal("expected invalid output error")
	}
	if err := (Config{BaseURL: "https://api.chartbrew.com", Token: "token", Output: "json"}).Validate(); err != nil {
		t.Fatalf("Validate() = %v", err)
	}
}
