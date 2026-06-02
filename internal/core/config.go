package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/aiwen/aw-cli/errs"
)

const (
	DefaultBaseURL  = "https://api.ipplus360.com"
	DefaultChannel  = "aw_cli"
	DefaultTimeout  = 10 * time.Second
	DefaultAccuracy = "city"
)

type CliConfig struct {
	BaseURL      string `json:"base_url"`
	APIKey       string `json:"api_key"`
	Channel      string `json:"channel"`
	Timeout      string `json:"timeout"`
	IPv4Accuracy string `json:"ipv4_accuracy"`
	IPv6Accuracy string `json:"ipv6_accuracy"`
}

type Overrides struct {
	ConfigPath string
	BaseURL    string
	APIKey     string
	Timeout    string
	Channel    string
}

func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aw-cli", "config.json"), nil
}

func Load(path string) (CliConfig, error) {
	if path == "" {
		defaultPath, err := DefaultConfigPath()
		if err != nil {
			return CliConfig{}, err
		}
		path = defaultPath
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return CliConfig{}, nil
	}
	if err != nil {
		return CliConfig{}, err
	}
	var cfg CliConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return CliConfig{}, errs.Config("failed to parse config file")
	}
	return cfg, nil
}

func Resolve(overrides Overrides) (CliConfig, error) {
	cfg, err := Load(overrides.ConfigPath)
	if err != nil {
		return CliConfig{}, err
	}

	cfg.BaseURL = firstNonEmpty(overrides.BaseURL, os.Getenv("AIWEN_API_BASE_URL"), cfg.BaseURL, DefaultBaseURL)
	cfg.APIKey = firstNonEmpty(overrides.APIKey, os.Getenv("AIWEN_API_KEY"), cfg.APIKey)
	cfg.Channel = firstNonEmpty(overrides.Channel, os.Getenv("AIWEN_CHANNEL"), cfg.Channel, DefaultChannel)
	cfg.Timeout = firstNonEmpty(overrides.Timeout, os.Getenv("AIWEN_TIMEOUT"), cfg.Timeout, DefaultTimeout.String())
	cfg.IPv4Accuracy = firstNonEmpty(os.Getenv("IPV4_ACCURACY"), cfg.IPv4Accuracy, DefaultAccuracy)
	cfg.IPv6Accuracy = firstNonEmpty(os.Getenv("IPV6_ACCURACY"), cfg.IPv6Accuracy, DefaultAccuracy)

	if !validAccuracy(cfg.IPv4Accuracy) {
		return CliConfig{}, errs.Config("invalid IPV4_ACCURACY; valid options are city, district, street")
	}
	if !validAccuracy(cfg.IPv6Accuracy) {
		return CliConfig{}, errs.Config("invalid IPV6_ACCURACY; valid options are city, district, street")
	}
	if _, err := time.ParseDuration(cfg.Timeout); err != nil {
		return CliConfig{}, errs.Config("invalid timeout duration")
	}
	return cfg, nil
}

func Write(path string, cfg CliConfig) error {
	if path == "" {
		defaultPath, err := DefaultConfigPath()
		if err != nil {
			return err
		}
		path = defaultPath
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}

func ParseTimeout(value string) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		return DefaultTimeout
	}
	return d
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func validAccuracy(value string) bool {
	return value == "city" || value == "district" || value == "street"
}
