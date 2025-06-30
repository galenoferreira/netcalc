// cidrcalc_test.go
package main

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestNetmask(t *testing.T) {
	// Para /24, a máscara deve ser 255.255.255.0
	want := binary.BigEndian.Uint32(net.ParseIP("255.255.255.0").To4())
	got := Netmask(24)
	if got != want {
		t.Errorf("Netmask(24) = %d; want %d", got, want)
	}
}

func TestIPConversions(t *testing.T) {
	ipStr := "10.1.2.3"
	u := IPToUint32(ipStr)
	back := Uint32ToIP(u)
	if back != ipStr {
		t.Errorf("IPToUint32/Uint32ToIP round-trip: got %q; want %q", back, ipStr)
	}
}
