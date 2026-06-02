package endpoint

import (
	"fmt"

	"github.com/aiwen/aw-cli/errs"
)

const (
	ActionLoc      = "loc"
	ActionCurrent  = "current"
	ActionScene    = "scene"
	ActionWhois    = "whois"
	ActionASN      = "asn"
	ActionHost     = "host"
	ActionRisk     = "risk"
	ActionIdentity = "identity"
	ActionIndustry = "industry"
	ActionAll      = "all"
)

type ActionSpec struct {
	Name         string
	SupportsIPv4 bool
	SupportsIPv6 bool
	AccuracyPath map[string]string
	IPv4Path     string
	IPv6Path     string
}

var Registry = map[string]ActionSpec{
	ActionLoc: {
		Name:         ActionLoc,
		SupportsIPv4: true,
		SupportsIPv6: true,
		AccuracyPath: map[string]string{
			"ipv4:city":     "/ip/geo/v1/city/",
			"ipv4:district": "/ip/geo/v1/district/",
			"ipv4:street":   "/ip/geo/v1/street/psi/",
			"ipv6:city":     "/ip/geo/v1/ipv6/",
			"ipv6:district": "/ip/geo/v1/ipv6/district/",
			"ipv6:street":   "/ip/geo/v1/ipv6/street/biz/",
		},
	},
	ActionScene: {
		Name:         ActionScene,
		SupportsIPv4: true,
		SupportsIPv6: true,
		IPv4Path:     "/ip/info/v1/scene/",
		IPv6Path:     "/ip/info/v1/ipv6Scene/",
	},
	ActionWhois: {
		Name:         ActionWhois,
		SupportsIPv4: true,
		IPv4Path:     "/ip/info/v1/ipWhois",
	},
	ActionASN: {
		Name:         ActionASN,
		SupportsIPv4: true,
		IPv4Path:     "/as/info/v1/asWhois",
	},
	ActionHost: {
		Name:         ActionHost,
		SupportsIPv4: true,
		IPv4Path:     "/ip/geo/v1/host/",
	},
	ActionRisk: {
		Name:         ActionRisk,
		SupportsIPv4: true,
		IPv4Path:     "/ip/info/v3/portrait/",
	},
	ActionIdentity: {
		Name:         ActionIdentity,
		SupportsIPv4: true,
		IPv4Path:     "/ip/info/v1/person/",
	},
	ActionIndustry: {
		Name:         ActionIndustry,
		SupportsIPv4: true,
		IPv4Path:     "/ip/info/v1/industry/",
	},
}

func Actions() []string {
	return []string{ActionLoc, ActionScene, ActionWhois, ActionASN, ActionHost, ActionRisk, ActionIdentity, ActionIndustry}
}

func Lookup(action string) (ActionSpec, error) {
	spec, ok := Registry[action]
	if !ok {
		return ActionSpec{}, errs.Validationf("unsupported action: %s", action)
	}
	return spec, nil
}

func Path(action, version, accuracy string) (string, error) {
	spec, err := Lookup(action)
	if err != nil {
		return "", err
	}
	if version == "ipv4" && !spec.SupportsIPv4 {
		return "", errs.Validationf("action %s does not support IPv4", action)
	}
	if version == "ipv6" && !spec.SupportsIPv6 {
		return "", errs.Validationf("action %s only supports IPv4", action)
	}
	if action == ActionLoc {
		if accuracy == "" {
			accuracy = "city"
		}
		path, ok := spec.AccuracyPath[fmt.Sprintf("%s:%s", version, accuracy)]
		if !ok {
			return "", errs.Validationf("invalid accuracy %q; valid options are city, district, street", accuracy)
		}
		return path, nil
	}
	if version == "ipv6" {
		return spec.IPv6Path, nil
	}
	return spec.IPv4Path, nil
}

func SupportsVersion(action, version string) error {
	spec, err := Lookup(action)
	if err != nil {
		return err
	}
	if version == "ipv6" && !spec.SupportsIPv6 {
		return errs.Validationf("action %s only supports IPv4", action)
	}
	if version == "ipv4" && !spec.SupportsIPv4 {
		return errs.Validationf("action %s does not support IPv4", action)
	}
	return nil
}
