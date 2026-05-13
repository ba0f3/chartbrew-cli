package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvBaseURL = "CHARTBREW_API_URL"
	EnvToken   = "CHARTBREW_TOKEN"
)

type Config struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
	Output  string `json:"output"`
}

type Sources struct {
	FlagBaseURL string
	FlagToken   string
	FlagOutput  string
	Env         map[string]string
	DotenvPath  string
	ConfigPath  string
}

func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".config", "chartbrew", "config.json")
	}
	return filepath.Join(home, ".config", "chartbrew", "config.json")
}

func Resolve(src Sources) (Config, error) {
	cfg := Config{Output: "json"}

	if src.ConfigPath != "" {
		fileCfg, err := readConfigFile(expandHome(src.ConfigPath))
		if err != nil {
			return Config{}, err
		}
		mergeConfig(&cfg, fileCfg)
	}

	if src.DotenvPath != "" {
		dotenv, err := readDotenv(src.DotenvPath)
		if err != nil {
			return Config{}, err
		}
		applyEnv(&cfg, dotenv)
	}

	env := src.Env
	if env == nil {
		env = osEnv()
	}
	applyEnv(&cfg, env)

	if src.FlagBaseURL != "" {
		cfg.BaseURL = src.FlagBaseURL
	}
	if src.FlagToken != "" {
		cfg.Token = src.FlagToken
	}
	if src.FlagOutput != "" {
		cfg.Output = src.FlagOutput
	}

	cfg.BaseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	cfg.Token = strings.TrimSpace(cfg.Token)
	cfg.Output = strings.TrimSpace(cfg.Output)
	return cfg, nil
}

func (c Config) Validate() error {
	if c.BaseURL == "" {
		return errors.New("missing Chartbrew API base URL; set --base-url or CHARTBREW_API_URL")
	}
	if c.Token == "" {
		return errors.New("missing Chartbrew API token; set --token or CHARTBREW_TOKEN")
	}
	switch c.Output {
	case "", "json", "markdown", "raw":
		return nil
	default:
		return fmt.Errorf("unsupported output format %q", c.Output)
	}
}

func readConfigFile(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return cfg, nil
	}
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return cfg, fmt.Errorf("parse config file: %w", err)
	}
	cfg.BaseURL = first(raw["base_url"], raw["api_url"])
	cfg.Token = raw["token"]
	cfg.Output = raw["output"]
	return cfg, nil
}

func readDotenv(path string) (map[string]string, error) {
	values := map[string]string{}
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return values, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan .env: %w", err)
	}
	return values, nil
}

func applyEnv(cfg *Config, env map[string]string) {
	if env[EnvBaseURL] != "" {
		cfg.BaseURL = env[EnvBaseURL]
	}
	if env[EnvToken] != "" {
		cfg.Token = env[EnvToken]
	}
}

func mergeConfig(dst *Config, src Config) {
	if src.BaseURL != "" {
		dst.BaseURL = src.BaseURL
	}
	if src.Token != "" {
		dst.Token = src.Token
	}
	if src.Output != "" {
		dst.Output = src.Output
	}
}

func osEnv() map[string]string {
	env := map[string]string{}
	for _, pair := range os.Environ() {
		key, value, ok := strings.Cut(pair, "=")
		if ok {
			env[key] = value
		}
	}
	return env
}

func expandHome(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return home
		}
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func first(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
