package batch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTextFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ips.txt")
	content := "8.8.8.8\n1.1.1.1\n# comment\n\n2001:4860:4860::8888\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ips, err := ParseFile(path, "ip")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(ips) != 3 {
		t.Errorf("expected 3 IPs, got %d: %v", len(ips), ips)
	}
	if ips[0] != "8.8.8.8" {
		t.Errorf("expected first IP 8.8.8.8, got %s", ips[0])
	}
}

func TestParseCSVFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ips.csv")
	content := "ip,name\n8.8.8.8,google\n1.1.1.1,cloudflare\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ips, err := ParseFile(path, "ip")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(ips) != 2 {
		t.Errorf("expected 2 IPs, got %d: %v", len(ips), ips)
	}
	if ips[0] != "8.8.8.8" {
		t.Errorf("expected first IP 8.8.8.8, got %s", ips[0])
	}
}

func TestParseCSVFileCustomColumn(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.csv")
	content := "address,port\n8.8.8.8,443\n1.1.1.1,80\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ips, err := ParseFile(path, "address")
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	if len(ips) != 2 {
		t.Errorf("expected 2 IPs, got %d", len(ips))
	}
	if ips[0] != "8.8.8.8" {
		t.Errorf("expected first IP 8.8.8.8, got %s", ips[0])
	}
}

func TestParseNonexistentFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.txt", "ip")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestFilterSpecial(t *testing.T) {
	ips := []string{"8.8.8.8", "192.168.1.1", "1.1.1.1", "127.0.0.1"}
	filtered := FilterSpecial(ips)
	if len(filtered) != 2 {
		t.Errorf("expected 2 public IPs, got %d: %v", len(filtered), filtered)
	}
	if filtered[0] != "8.8.8.8" || filtered[1] != "1.1.1.1" {
		t.Errorf("expected [8.8.8.8, 1.1.1.1], got %v", filtered)
	}
}

func TestFilterSpecialEmpty(t *testing.T) {
	ips := []string{}
	filtered := FilterSpecial(ips)
	if len(filtered) != 0 {
		t.Errorf("expected 0 IPs, got %d", len(filtered))
	}
}

func TestFilterSpecialInvalidIPs(t *testing.T) {
	ips := []string{"not-an-ip", "8.8.8.8"}
	filtered := FilterSpecial(ips)
	if len(filtered) != 2 {
		t.Errorf("expected invalid IPs to pass through, got %d: %v", len(filtered), filtered)
	}
}
