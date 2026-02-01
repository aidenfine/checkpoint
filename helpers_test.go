package checkpointmiddleware_test

import (
	"net"
	"testing"

	checkpointmiddleware "github.com/aidenfine/checkpoint"
)

func TestCanonicalizeWithIpv4(t *testing.T) {
	ipv4 := "185.214.81.61"

	masked_ipv4 := checkpointmiddleware.CanonicalizeIPFunc(ipv4)
	if masked_ipv4 != ipv4 {
		t.Errorf("ipv4 do not match, originalIpv4: %s masked_ipv4: %s", ipv4, masked_ipv4)
	}

	ipv6 := "d591:29ff:c1c5:045d:21a5:fad9:d322:5459"
	ipv6_parsed := net.ParseIP(ipv6)

	masked_ipv6 := checkpointmiddleware.CanonicalizeIPFunc(ipv6)
	if ipv6_parsed.Mask(net.CIDRMask(64, 128)).String() != masked_ipv6 {
		t.Errorf("ipv6 do not match, originalIpv6: %s masked_ipv6: %s", ipv4, masked_ipv4)

	}

}
