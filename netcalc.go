// Package main implements netcalc, a CLI IPv4 subnet calculator with
// support for interactive menus, single-command invocation, and detailed output.
/*...
 * netcalc - IPv4 Subnet Calculator...
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
 *
 */

package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// printVersionInfo outputs the build version details.
func printVersionInfo() {
	fmt.Printf("Build Time : %s\nGit Commit : %s\nGit Branch : %s\n", buildTime, gitCommit, gitBranch)
}

// Version information populated via ldflags at build time.
var buildTime string // build timestamp
var gitCommit string // git commit hash
var gitBranch string // git branch name

// ErrHelp is returned by parseInput when the user requests help.
var ErrHelp = errors.New("help requested")

// parseInput parses the provided command-line arguments and returns
// an IP string and prefix length. Supports:
//   - No arguments: returns ErrHelp
//   - Single argument with CIDR format "IP/prefix"
//   - Two arguments: "IP mask" (dotted mask) or "IP prefix" (numeric)
//
// It validates IP format and prefix/mask ranges.
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
				if v < 0 || v > 255 {
					return "", 0, fmt.Errorf("mask octet out of range: %d", v)
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

// NetCalc computes and displays comprehensive subnet details:
//   - Network and broadcast addresses
//   - Usable host range and total hosts
//   - Wildcard mask and private-range overflow warnings
//
// The output is formatted with emojis and aligned columns.
func NetCalc(ipStr string, bits uint, showBin, showHex bool) {
	// Convert values
	maskUint := Netmask(bits)
	ipUint := IPToUint32(ipStr)
	network := ipUint & maskUint
	broadcast := network | ^maskUint
	firstIP := network + 1
	lastIP := broadcast - 1
	totalHosts := (1 << (32 - bits)) - 2
	wildcard := ^maskUint

	// privateRanges defines RFC1918 private IPv4 blocks to check for overflow.
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

	// ANSI color codes
	const (
		colorReset  = "\033[0m"
		colorYellow = "\033[33m"
		colorCyan   = "\033[36m"
		colorGreen  = "\033[32m"
	)

	// Input section
	fmt.Printf("%s❯ netcalc %s/%d%s\n\n", colorYellow, ipStr, bits, colorReset)
	fmt.Println("---[ 🔍 INPUT ]----------------------------------")
	fmt.Printf("Address               : %s%s/%d%s\n\n", colorYellow, ipStr, bits, colorReset)

	// Network section
	fmt.Println("---[ 🌐 NETWORK ]----------------------------------")
	fmt.Printf("Network               : %s%s%s\n", colorCyan, Uint32ToIP(network), colorReset)
	fmt.Printf("Broadcast             : %s%s%s\n", colorCyan, Uint32ToIP(broadcast), colorReset)
	fmt.Printf("Usable Hosts          : %s%s - %s%s\n", colorCyan, Uint32ToIP(firstIP), Uint32ToIP(lastIP), colorReset)
	fmt.Printf("Total Hosts           : %s%d%s\n\n", colorGreen, totalHosts, colorReset)

	// Mask section
	fmt.Println("---[ 🔢 MASK ]--------------------------------------")
	fmt.Printf("Network Mask          : %s%s%s\n", colorYellow, Uint32ToIP(maskUint), colorReset)
	fmt.Printf("Wildcard Mask         : %s%s%s\n", colorYellow, Uint32ToIP(wildcard), colorReset)
	fmt.Printf("Hexadecimal Mask      : %s0x%X%s\n", colorYellow, maskUint, colorReset)

	// Optional binary display
	if showBin {
		fmt.Println("\n---[ 🔢 BINARY ]---------------------------------")
		fmt.Printf("IP (binary)           : %032b\n", ipUint)
		fmt.Printf("Netmask (binary)      : %032b\n", maskUint)
		fmt.Printf("Network (binary)      : %032b\n", network)
		fmt.Printf("Broadcast (binary)    : %032b\n", broadcast)
	}

	// Optional hexadecimal display
	if showHex {
		fmt.Println("\n---[ 🔢 HEXADECIMAL ]-------------------------------")
		fmt.Printf("IP (hex)              : 0x%X\n", ipUint)
		fmt.Printf("Netmask (hex)         : 0x%X\n", maskUint)
		fmt.Printf("Network (hex)         : 0x%X\n", network)
		fmt.Printf("Broadcast (hex)       : 0x%X\n", broadcast)
	}
}

// Netmask returns a 32-bit network mask for a given CIDR prefix length.
// It shifts a 32-bit all-ones mask left by (32 - bits).
func Netmask(bits uint) uint32 {
	return ^uint32(0) << (32 - bits)
}

// IPToUint32 converts an IPv4 string (dotted quad) to its 32-bit uint representation.
func IPToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr).To4()
	return binary.BigEndian.Uint32(ip)
}

// Uint32ToIP converts a 32-bit integer to an IPv4 dotted-quad string.
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

// printManual outputs the embedded manual text to stdout
func printManual() {
	fmt.Print(manualText)
}

// main is the CLI entry point. It handles:
//   - --help and no-argument cases to show the manual
//   - Parsing input arguments
//   - Invoking NetCalc for subnet calculation
func main() {
	// Define CLI flags
	showBin := flag.Bool("b", false, "Show addresses in binary")
	showHex := flag.Bool("h", false, "Show addresses in hexadecimal")
	showVersion := flag.Bool("version", false, "Display version information")
	flag.Parse()

	if *showVersion {
		printVersionInfo()
		return
	}

	// Remaining args after flags
	args := flag.Args()

	// If the user provided no args or asked for help, show the manual.
	if len(args) == 0 || (len(args) > 0 && args[0] == "--help") {
		printManual()
		return
	}

	// Otherwise, attempt to parse as CIDR or IP/mask.
	ip, prefix, err := parseInput(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	NetCalc(ip, prefix, *showBin, *showHex)
}
