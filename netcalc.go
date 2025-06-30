/*
 * netcalc - IPv4 Subnet Calculator..
 *
 * Version Control:
 *   Build Time : ${BUILD_TIME}
 *   Git Commit : ${GIT_COMMIT}
 *   Git Branch : ${GIT_BRANCH}
 *
 * These variables can be set at build time using ldflags, for example:
 *   go build -ldflags "-X main.buildTime=$(date +%Y-%m-%dT%H:%M:%SZ) \
 *                       -X main.gitCommit=$(git rev-parse --short HEAD) \
 *                       -X main.gitBranch=$(git rev-parse --abbrev-ref HEAD)" \
 *     -o netcalc
 */

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var buildTime string // build timestamp
var gitCommit string // git commit hash
var gitBranch string // git branch name

var ErrHelp = errors.New("help requested")

// parseInput parses command-line arguments into an IP string and prefix length.
func parseInput(args []string) (string, uint, error) {
	if len(args) == 0 {
		return "", 0, ErrHelp
	}
	// Single argument with slash notation
	if len(args) == 1 && strings.Contains(args[0], "/") {
		parts := strings.SplitN(args[0], "/", 2)
		ip := parts[0]
		prefix, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, fmt.Errorf("invalid prefix: %s", parts[1])
		}
		return ip, uint(prefix), nil
	}
	// Two arguments: IP and mask or CIDR
	if len(args) == 2 {
		ip := args[0]
		// Dotted mask
		if strings.Contains(args[1], ".") {
			// Convert mask to prefix
			maskParts := strings.Split(args[1], ".")
			if len(maskParts) != 4 {
				return "", 0, fmt.Errorf("invalid mask: %s", args[1])
			}
			var maskBytes [4]byte
			for i, p := range maskParts {
				v, err := strconv.Atoi(p)
				if err != nil {
					return "", 0, fmt.Errorf("invalid mask octet: %s", p)
				}
				maskBytes[i] = byte(v)
			}
			ones, bits := net.IPMask(maskBytes[:]).Size()
			if bits != 32 {
				return "", 0, fmt.Errorf("invalid mask size: %d", bits)
			}
			return ip, uint(ones), nil
		}
		// Numeric CIDR
		prefix, err := strconv.Atoi(args[1])
		if err != nil {
			return "", 0, fmt.Errorf("invalid prefix: %s", args[1])
		}
		return ip, uint(prefix), nil
	}
	return "", 0, fmt.Errorf("unexpected arguments: %v", args)
}

// netcalc computes and displays detailed subnet information for the given IP and prefix length.
// It calculates the network address, broadcast address, usable host range, wildcard mask, and total hosts.
func NetCalc(ipStr string, bits uint) {
	// Convert values
	maskUint := Netmask(bits)
	ipUint := IPToUint32(ipStr)
	network := ipUint & maskUint
	broadcast := network | ^maskUint
	firstIP := network + 1
	lastIP := broadcast - 1
	totalHosts := (1 << (32 - bits)) - 2
	wildcard := ^maskUint

	// Define RFC1918 private network ranges
	privateRanges := []struct {
		base   uint32 // network base
		prefix uint   // CIDR prefix length
	}{
		{IPToUint32("10.0.0.0"), 8},
		{IPToUint32("172.16.0.0"), 12},
		{IPToUint32("192.168.0.0"), 16},
	}

	for _, r := range privateRanges {
		// check if input IP is within this private block
		blockStart := r.base
		blockEnd := r.base | ^Netmask(r.prefix)
		if ipUint >= blockStart && ipUint <= blockEnd {
			// warn if the computed subnet lies outside the private block
			if network < blockStart || broadcast > blockEnd {
				fmt.Printf("🚨 Warning: Subnet %s/%d extends outside private block %s/%d\n",
					Uint32ToIP(network), bits,
					Uint32ToIP(blockStart), r.prefix)
			}
			break
		}
	}

	// Detailed, emoji-rich output
	fmt.Printf("%-25s %s\n", "🔍 Input IP:", ipStr)
	fmt.Printf("%-25s %s (/ %d)\n", "🔢 Netmask:", Uint32ToIP(maskUint), bits)
	fmt.Printf("%-25s 0x%X\n", "🛡️ Mask (hex):", maskUint)
	fmt.Printf("%-25s %s\n", "✖ Wildcard Mask:", Uint32ToIP(wildcard))
	fmt.Printf("%-25s %s/%d\n", "🌐 Network Address:", Uint32ToIP(network), bits)
	fmt.Printf("%-25s %s\n", "📡 Broadcast Address:", Uint32ToIP(broadcast))
	fmt.Printf("%-25s %s - %s\n", "↕️ Usable Host Range:", Uint32ToIP(firstIP), Uint32ToIP(lastIP))
	fmt.Printf("%-25s %d\n", "📊 Total Usable Hosts:", totalHosts)
}

// Netmask returns a 32-bit network mask for the given prefix length.
func Netmask(bits uint) uint32 {
	return ^uint32(0) << (32 - bits)
}

// IPToUint32 converts a dotted IP string to a 32-bit integer.
func IPToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr).To4()
	return binary.BigEndian.Uint32(ip)
}

// Uint32ToIP converts a 32-bit integer to a dotted IP string.
func Uint32ToIP(n uint32) string {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], n)
	return net.IP(b[:]).String()
}

const manualText = `
NETCALC(1)   User Commands   NETCALC(1)

NAME
    netcalc — IPv4 subnet calculator CLI tool

SYNOPSIS
    netcalc [options] <IP address>/<prefix length>

DESCRIPTION
    netcalc is a command-line tool that computes and displays detailed subnet information
    for a given IPv4 address and prefix length.

EXAMPLES
    netcalc 192.168.1.10/24
    netcalc 10.0.0.1/8

`

// printManual prints the full built-in manual to stdout.
func printManual() {
	fmt.Print(manualText)
}

// main is the entry point for the netcalc tool.
// It displays the manual when run without arguments or with help flags.
func main() {
	// If the user provided no args or asked for help, show the manual.
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printManual()
		return
	}

	// Otherwise, attempt to parse as CIDR or IP/mask.
	ip, prefix, err := parseInput(os.Args[1:])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	NetCalc(ip, prefix)
}
