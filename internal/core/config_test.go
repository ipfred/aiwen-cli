package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	path, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("DefaultConfigPath error: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty config path")
	}
}

func TestLoadNonexistent(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("Load nonexistent file error: %v", err)
	}
	if cfg.BaseURL != "" {
		t.Error("expected empty config for nonexistent file")
	}
}

func TestWriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	cfg := CliConfig{
		BaseURL:      "https://api.example.com",
		APIKey:       "test-key",
		Channel:      "test_channel",
		Timeout:      "5s",
		IPv4Accuracy: "district",
		IPv6Accuracy: "city",
	}
	if err := Write(path, cfg); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if loaded.BaseURL != cfg.BaseURL {
		t.Errorf("BaseURL mismatch: got %q, want %q", loaded.BaseURL, cfg.BaseURL)
	}
	if loaded.APIKey != cfg.APIKey {
		t.Errorf("APIKey mismatch: got %q, want %q", loaded.APIKey, cfg.APIKey)
	}
	if loaded.IPv4Accuracy != cfg.IPv4Accuracy {
		t.Errorf("IPv4Accuracy mismatch: got %q, want %q", loaded.IPv4Accuracy, cfg.IPv4Accuracy)
	}
}

func TestResolveDefaults(t *testing.T) {
	cfg, err := Resolve(Overrides{})
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("expected default BaseURL %q, got %q", DefaultBaseURL, cfg.BaseURL)
	}
	if cfg.Channel != DefaultChannel {
		t.Errorf("expected default Channel %q, got %q", DefaultChannel, cfg.Channel)
	}
	if cfg.IPv4Accuracy != DefaultAccuracy {
		t.Errorf("expected default IPv4Accuracy %q, got %q", DefaultAccuracy, cfg.IPv4Accuracy)
	}
}

func TestResolveWithOverrides(t *testing.T) {
	cfg, err := Resolve(Overrides{
		BaseURL: "https://custom.api.com",
		APIKey:  "override-key",
		Channel: "custom",
	})
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if cfg.BaseURL != "https://custom.api.com" {
		t.Errorf("expected override BaseURL, got %q", cfg.BaseURL)
	}
	if cfg.APIKey != "override-key" {
		t.Errorf("expected override APIKey, got %q", cfg.APIKey)
	}
	if cfg.Channel != "custom" {
		t.Errorf("expected override Channel, got %q", cfg.Channel)
	}
}

func TestResolveWithEnvOverrides(t *testing.T) {
	os.Setenv("AIWEN_API_KEY", "env-key")
	os.Setenv("AIWEN_API_BASE_URL", "https://env.api.com")
	defer os.Unsetenv("AIWEN_API_KEY")
	defer os.Unsetenv("AIWEN_API_BASE_URL")

	cfg, err := Resolve(Overrides{})
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if cfg.APIKey != "env-key" {
		t.Errorf("expected env APIKey, got %q", cfg.APIKey)
	}
	if cfg.BaseURL != "https://env.api.com" {
		t.Errorf("expected env BaseURL, got %q", cfg.BaseURL)
	}
}

func TestResolveInvalidAccuracy(t *testing.T) {
	_, err := Resolve(Overrides{})
	os.Setenv("IPV4_ACCURACY", "invalid")
	defer os.Unsetenv("IPV4_ACCURACY")
	_, err = Resolve(Overrides{})
	if err == nil {
		t.Error("expected error for invalid accuracy")
	}
}

func TestRedactSecret(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"short", "abc", "***"},
		{"long", "abcdefgh", "abc***fgh"},
		{"6char", "abcdef", "***"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactSecret(tt.input); got != tt.want {
				t.Errorf("RedactSecret(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseTimeout(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid", "5s", "5s"},
		{"valid minutes", "2m", "2m0s"},
		{"invalid", "bad", "10s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimeout(tt.input)
			if got.String() != tt.want {
				t.Errorf("ParseTimeout(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}
