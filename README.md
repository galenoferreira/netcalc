# 🌐 netcalc — IPv4 Subnet Calculator

A command-line tool to perform comprehensive IPv4 subnet calculations and network operations.

## 🚀 Features

- Calculate network mask, network address, broadcast address, usable host range, wildcard mask, and total hosts.
- Interactive menu when run without arguments.
- Shortcuts for non-interactive use:
  - `-c <ip>/<cidr>` — Calculate via CIDR notation.
  - `-i <ip> <mask>` — Calculate via IP and mask.
  - `-w <ip>/<cidr>` — Show wildcard mask.
  - `-W <wildcard>` — Convert wildcard mask to CIDR.
  - `-I <network> <ip>` — Check IP inclusion in network.
  - `-l <network>` — List all usable hosts.
  - `-n <network>` — Show previous and next subnets.
  - `-b <ip>/<cidr>` — Display IP and mask in binary.
  - `-r <network>` — Print reverse DNS zone.
  - `--help` / `-h` — Show brief usage.
  - `--manual` / `-M` — Show full embedded manual.
  - `--version` — Display build version, commit, and branch.

## 📦 Installation

Build from source (requires Go):

```bash
git clone https://github.com/galenoferreira/netcalc.git
cd netcalc
go build -ldflags "-X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
                   -X main.gitCommit=$(git rev-parse --short HEAD) \
                   -X main.gitBranch=$(git rev-parse --abbrev-ref HEAD)" \
         -o netcalc netcalc.go
```

Or use the provided `build.sh` script:

```bash
chmod +x build.sh
./build.sh
```

## 🛠 Usage

### Interactive mode

```bash
./netcalc
```

### Non-interactive shortcuts

```bash
./netcalc -c 192.168.0.1/24
./netcalc -i 192.168.0.1 255.255.255.0
```

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## 🤝 Contributing

Please open issues and pull requests on GitHub:  
https://github.com/galenoferreira/netcalc


## 🎖️ Credits

Created and maintained by Galeno Garbe.
