package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const helpText = `CIDRCalc – Usage Guide

Usage:
  cidrcalc <ip>/<cidr>
  cidrcalc <ip> <mask>

Examples:
  cidrcalc 192.168.0.1/24
  cidrcalc 10.0.0.5 255.255.255.0

Options:
  -h, --help    Show this help message

`

// Netmask calculates the network mask based on CIDR
func Netmask(bits uint) uint32 {
	return ^(uint32(0xFFFFFFFF) >> bits)
}

// Uint32ToIP converts a uint32 to an IP string
func Uint32ToIP(ipUint uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipUint)
	return ip.String()
}

// IPToUint32 converts an IP string to a uint32
func IPToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr).To4()
	return binary.BigEndian.Uint32(ip)
}

// maskToPrefix converts a dotted mask (e.g. "255.255.255.0") to a CIDR prefix length.
func maskToPrefix(mask string) (uint, error) {
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid mask format")
	}
	var maskBytes [4]byte
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || v < 0 || v > 255 {
			return 0, fmt.Errorf("invalid mask octet: %s", p)
		}
		maskBytes[i] = byte(v)
	}
	ones, bits := net.IPMask(maskBytes[:]).Size()
	if bits != 32 {
		return 0, fmt.Errorf("unexpected mask size %d", bits)
	}
	return uint(ones), nil
}

// CIDRCalc is the main calculation function
func CIDRCalc(ipStr string, bits uint) {
	ip := IPToUint32(ipStr)
	mask := Netmask(bits)

	network := ip & mask
	broadcast := network | ^mask
	firstIP := network + 1
	lastIP := broadcast - 1
	totalHosts := (1 << (32 - bits)) - 2

	fmt.Println("Input IP:", ipStr)
	fmt.Println("Netmask:", Uint32ToIP(mask))
	fmt.Println("Network Address:", Uint32ToIP(network))
	fmt.Println("First Usable IP:", Uint32ToIP(firstIP))
	fmt.Println("Last Usable IP:", Uint32ToIP(lastIP))
	fmt.Println("Broadcast Address:", Uint32ToIP(broadcast))
	fmt.Println("Total Valid Hosts:", totalHosts)
}

func main() {
	// Expect either:
	//   cidrcalc 192.168.0.1/24
	//   cidrcalc 192.168.0.1 255.255.255.0
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(helpText)
		os.Exit(0)
	}
	var ip string
	var prefix uint
	if len(args) == 1 && strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		ip = parts[0]
		v, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Invalid CIDR prefix:", parts[1])
			os.Exit(1)
		}
		prefix = uint(v)
	} else if len(args) == 2 {
		ip = args[0]
		if strings.Contains(args[1], ".") {
			// dotted mask
			v, err := maskToPrefix(args[1])
			if err != nil {
				fmt.Println("Invalid mask:", err)
				os.Exit(1)
			}
			prefix = v
		} else {
			// numeric CIDR
			v, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid CIDR prefix:", args[1])
				os.Exit(1)
			}
			prefix = uint(v)
		}
	} else {
		fmt.Println("Usage: cidrcalc <ip>/<cidr> OR cidrcalc <ip> <mask>")
		os.Exit(1)
	}
	CIDRCalc(ip, prefix)
}
