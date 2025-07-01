# 🌐 netcalc — IPv4 Subnet Calculator

[![Build Status](https://github.com/galenoferreira/cidrcalc/actions/workflows/ci.yml/badge.svg)](https://github.com/galenoferreira/cidrcalc/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/galenoferreira/cidrcalc)](https://goreportcard.com/report/github.com/galenoferreira/cidrcalc)
[![codecov](https://codecov.io/gh/galenoferreira/cidrcalc/branch/master/graph/badge.svg?token=)](https://codecov.io/gh/galenoferreira/cidrcalc)

A single-binary CLI tool for comprehensive IPv4 subnet calculations.

## 🚀 Features

- **Subnet Calculations**: network, broadcast, usable host range, total hosts.
- **Mask Utilities**: network mask, wildcard mask, hexadecimal mask.
- **Optional Displays**:
  - `-b` show binary representation of IP, mask, network, and broadcast.
  - `-h` show hexadecimal representation of IP, mask, network, and broadcast.
- **Interactive Mode**: run without arguments for a menu-driven interface.
- **Flags** for quick operations:
  - `-c <IP/CIDR>` — CIDR notation calculation.
  - `-i <IP> <mask>` — IP + dotted mask or numeric prefix.
  - `-w <IP/CIDR>` — wildcard mask only.
  - `-W <wildcard>` — convert wildcard mask to prefix.
  - `-I <network> <IP>` — check IP inclusion in network.
  - `-l <network>` — list all usable hosts.
  - `-n <network>` — previous and next subnet.
  - `--help`, `-h` — display short usage.
  - `--manual`, `-M` — display full manual.
  - `--version` — show build time, commit, and branch.
  - `--version` — show build time, git commit, and branch information.

## 📦 Installation

Build from source (requires Go):

```bash
go build -ldflags "-X main.buildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
                   -X main.gitCommit=$(git rev-parse --short HEAD) \
                   -X main.gitBranch=$(git rev-parse --abbrev-ref HEAD')" \
                   -o netcalc netcalc.go

# Display version information
./netcalc --version
```
