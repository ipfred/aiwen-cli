package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"json", FormatJSON, false},
		{"ndjson", FormatNDJSON, false},
		{"table", FormatTable, false},
		{"csv", FormatCSV, false},
		{"invalid", "xml", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.format)
			if tt.wantErr && err == nil {
				t.Errorf("ValidateFormat(%q) expected error", tt.format)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidateFormat(%q) unexpected error: %v", tt.format, err)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]any{"country": "US", "city": "Mountain View"}
	err := Write(&buf, data, FormatJSON)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if result["country"] != "US" {
		t.Errorf("expected country=US, got %v", result["country"])
	}
}

func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]any{"country": "US", "city": "Mountain View"}
	err := Write(&buf, data, FormatTable)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "country") || !strings.Contains(output, "US") {
		t.Errorf("table output missing expected content: %s", output)
	}
}

func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]any{"country": "US", "city": "Mountain View"}
	err := Write(&buf, data, FormatCSV)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "country") || !strings.Contains(output, "US") {
		t.Errorf("CSV output missing expected content: %s", output)
	}
}

func TestWriteNDJSON(t *testing.T) {
	var buf bytes.Buffer
	data := []any{
		map[string]any{"ip": "8.8.8.8"},
		map[string]any{"ip": "1.1.1.1"},
	}
	err := Write(&buf, data, FormatNDJSON)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestFilter(t *testing.T) {
	data := map[string]any{
		"ok":   true,
		"data": map[string]any{"country": "US"},
	}
	result := Filter(data, ".data.country")
	if result != "US" {
		t.Errorf("Filter(.data.country) = %v, want US", result)
	}
}

func TestFilterEmpty(t *testing.T) {
	data := map[string]any{"country": "US"}
	result := Filter(data, "")
	if result == nil {
		t.Error("Filter with empty expr should return data")
	}
}

func TestEnvelopeJSON(t *testing.T) {
	var buf bytes.Buffer
	env := Envelope{
		OK:     true,
		Action: "loc",
		IP:     "8.8.8.8",
		Data:   map[string]any{"country": "US"},
	}
	err := Write(&buf, env, FormatJSON)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if result["ok"] != true {
		t.Error("expected ok=true")
	}
	if result["action"] != "loc" {
		t.Error("expected action=loc")
	}
}
