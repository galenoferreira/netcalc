package main

import "testing"

func TestNetmask(t *testing.T) {
	got := Netmask(24)
	want := uint32(0xFFFFFF00)
	if got != want {
		t.Errorf("Netmask(24) = %X; want %X", got, want)
	}
}
