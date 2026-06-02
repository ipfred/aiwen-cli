package endpoint

import (
	"testing"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{"loc action", ActionLoc, false},
		{"scene action", ActionScene, false},
		{"whois action", ActionWhois, false},
		{"asn action", ActionASN, false},
		{"host action", ActionHost, false},
		{"risk action", ActionRisk, false},
		{"identity action", ActionIdentity, false},
		{"industry action", ActionIndustry, false},
		{"invalid action", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Lookup(tt.action)
			if tt.wantErr && err == nil {
				t.Errorf("Lookup(%q) expected error, got nil", tt.action)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Lookup(%q) unexpected error: %v", tt.action, err)
			}
		})
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name     string
		action   string
		version  string
		accuracy string
		want     string
		wantErr  bool
	}{
		{"loc ipv4 city", ActionLoc, "ipv4", "city", "/ip/geo/v1/city/", false},
		{"loc ipv4 district", ActionLoc, "ipv4", "district", "/ip/geo/v1/district/", false},
		{"loc ipv4 street", ActionLoc, "ipv4", "street", "/ip/geo/v1/street/psi/", false},
		{"loc ipv6 city", ActionLoc, "ipv6", "city", "/ip/geo/v1/ipv6/", false},
		{"loc ipv6 district", ActionLoc, "ipv6", "district", "/ip/geo/v1/ipv6/district/", false},
		{"loc ipv6 street", ActionLoc, "ipv6", "street", "/ip/geo/v1/ipv6/street/biz/", false},
		{"scene ipv4", ActionScene, "ipv4", "", "/ip/info/v1/scene/", false},
		{"scene ipv6", ActionScene, "ipv6", "", "/ip/info/v1/ipv6Scene/", false},
		{"whois ipv4", ActionWhois, "ipv4", "", "/ip/info/v1/ipWhois", false},
		{"asn ipv4", ActionASN, "ipv4", "", "/as/info/v1/asWhois", false},
		{"host ipv4", ActionHost, "ipv4", "", "/ip/geo/v1/host/", false},
		{"risk ipv4", ActionRisk, "ipv4", "", "/ip/info/v3/portrait/", false},
		{"identity ipv4", ActionIdentity, "ipv4", "", "/ip/info/v1/person/", false},
		{"industry ipv4", ActionIndustry, "ipv4", "", "/ip/info/v1/industry/", false},
		{"whois ipv6 not supported", ActionWhois, "ipv6", "", "", true},
		{"risk ipv6 not supported", ActionRisk, "ipv6", "", "", true},
		{"loc default accuracy", ActionLoc, "ipv4", "", "/ip/geo/v1/city/", false},
		{"loc invalid accuracy", ActionLoc, "ipv4", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Path(tt.action, tt.version, tt.accuracy)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Path(%q, %q, %q) expected error, got %q", tt.action, tt.version, tt.accuracy, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Path(%q, %q, %q) unexpected error: %v", tt.action, tt.version, tt.accuracy, err)
			}
			if got != tt.want {
				t.Errorf("Path(%q, %q, %q) = %q, want %q", tt.action, tt.version, tt.accuracy, got, tt.want)
			}
		})
	}
}

func TestSupportsVersion(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		version string
		wantErr bool
	}{
		{"loc supports ipv4", ActionLoc, "ipv4", false},
		{"loc supports ipv6", ActionLoc, "ipv6", false},
		{"risk supports ipv4", ActionRisk, "ipv4", false},
		{"risk not supports ipv6", ActionRisk, "ipv6", true},
		{"whois supports ipv4", ActionWhois, "ipv4", false},
		{"whois not supports ipv6", ActionWhois, "ipv6", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SupportsVersion(tt.action, tt.version)
			if tt.wantErr && err == nil {
				t.Errorf("SupportsVersion(%q, %q) expected error, got nil", tt.action, tt.version)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("SupportsVersion(%q, %q) unexpected error: %v", tt.action, tt.version, err)
			}
		})
	}
}

func TestActions(t *testing.T) {
	actions := Actions()
	if len(actions) != 8 {
		t.Errorf("Actions() returned %d actions, want 8", len(actions))
	}
}
