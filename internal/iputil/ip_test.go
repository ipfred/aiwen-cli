package iputil

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantIPv string
		wantErr bool
	}{
		{"valid IPv4", "8.8.8.8", "ipv4", false},
		{"valid IPv4 private", "192.168.1.1", "ipv4", false},
		{"valid IPv6", "2001:4860:4860::8888", "ipv6", false},
		{"valid IPv6 loopback", "::1", "ipv6", false},
		{"invalid string", "not-an-ip", "", true},
		{"empty string", "", "", true},
		{"hostname", "example.com", "", true},
		{"cidr", "8.8.8.0/24", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			ver, err := Version(tt.input)
			if err != nil {
				t.Fatalf("Version(%q) unexpected error: %v", tt.input, err)
			}
			if ver != tt.wantIPv {
				t.Errorf("Version(%q) = %q, want %q", tt.input, ver, tt.wantIPv)
			}
			if addr.Raw != tt.input {
				t.Errorf("Parse(%q).Raw = %q, want %q", tt.input, addr.Raw, tt.input)
			}
		})
	}
}

func TestIsSpecial(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"private", "192.168.1.1", true},
		{"loopback", "127.0.0.1", true},
		{"link-local", "169.254.1.1", true},
		{"multicast", "224.0.0.1", true},
		{"unspecified", "0.0.0.0", true},
		{"public", "8.8.8.8", false},
		{"public v6", "2001:4860:4860::8888", false},
		{"v6 loopback", "::1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if got := IsSpecial(addr.Addr); got != tt.want {
				t.Errorf("IsSpecial(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
