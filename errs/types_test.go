package errs

import (
	"testing"
)

func TestValidation(t *testing.T) {
	err := Validation("invalid input")
	if err.Code != ExitValidation {
		t.Errorf("expected exit code %d, got %d", ExitValidation, err.Code)
	}
	if err.Problem.Type != TypeValidation {
		t.Errorf("expected type %s, got %s", TypeValidation, err.Problem.Type)
	}
	if err.Problem.Message != "invalid input" {
		t.Errorf("expected message 'invalid input', got %s", err.Problem.Message)
	}
}

func TestValidationf(t *testing.T) {
	err := Validationf("invalid IP: %s", "bad")
	if err.Problem.Message != "invalid IP: bad" {
		t.Errorf("expected formatted message, got %s", err.Problem.Message)
	}
}

func TestConfig(t *testing.T) {
	err := Config("missing key")
	if err.Code != ExitConfig {
		t.Errorf("expected exit code %d, got %d", ExitConfig, err.Code)
	}
	if err.Problem.Type != TypeConfig {
		t.Errorf("expected type %s, got %s", TypeConfig, err.Problem.Type)
	}
}

func TestAPI(t *testing.T) {
	err := API("upstream error")
	if err.Code != ExitAPI {
		t.Errorf("expected exit code %d, got %d", ExitAPI, err.Code)
	}
}

func TestNetwork(t *testing.T) {
	err := Network("timeout")
	if err.Code != ExitNetwork {
		t.Errorf("expected exit code %d, got %d", ExitNetwork, err.Code)
	}
}

func TestParse(t *testing.T) {
	err := Parse("invalid json")
	if err.Code != ExitAPI {
		t.Errorf("expected exit code %d, got %d", ExitAPI, err.Code)
	}
	if err.Problem.Type != TypeParseError {
		t.Errorf("expected type %s, got %s", TypeParseError, err.Problem.Type)
	}
}

func TestWithDetail(t *testing.T) {
	err := WithDetail(ExitValidation, TypeValidation, "bad input", map[string]any{"field": "ip"})
	if err.Problem.Detail["field"] != "ip" {
		t.Errorf("expected detail field=ip, got %v", err.Problem.Detail)
	}
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name string
		code int
		want int
	}{
		{"internal", ExitInternal, 1},
		{"validation", ExitValidation, 2},
		{"config", ExitConfig, 3},
		{"api", ExitAPI, 4},
		{"network", ExitNetwork, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.want {
				t.Errorf("exit code = %d, want %d", tt.code, tt.want)
			}
		})
	}
}
