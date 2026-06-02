package iputil

import (
	"net/netip"

	"github.com/aiwen/aw-cli/errs"
)

type Addr struct {
	Raw  string
	Addr netip.Addr
}

func Parse(ip string) (Addr, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return Addr{}, errs.Validationf("invalid IP address: %s", ip)
	}
	return Addr{Raw: ip, Addr: addr}, nil
}

func Version(ip string) (string, error) {
	addr, err := Parse(ip)
	if err != nil {
		return "", err
	}
	if addr.Addr.Is4() {
		return "ipv4", nil
	}
	return "ipv6", nil
}

func IsSpecial(addr netip.Addr) bool {
	return addr.IsPrivate() ||
		addr.IsLoopback() ||
		addr.IsMulticast() ||
		addr.IsUnspecified() ||
		addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast()
}
