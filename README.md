# 🌐 netcalc — IPv4 Subnet Calculator by GG

[![CI/CD Pipeline](https://github.com/galenoferreira/netcalc/actions/workflows/main.yml/badge.svg)](https://github.com/galenoferreira/netcalc/actions/workflows/main.yml)

A single-binary CLI tool for comprehensive IPv4 subnet calculations.

## 🚀 Features

- **Subnet Calculations**: network, broadcast, usable host range, total hosts.
- **Mask Utilities**: network mask, wildcard mask, hexadecimal mask.

## Installation

Build from source (requires Go):

```bash
go build -ldflags "-X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
                   -X main.gitCommit=$(git rev-parse --short HEAD) \
                   -X main.gitBranch=$(git rev-parse --abbrev-ref HEAD)" \
           -o netcalc netcalc.go

# Display version information
./netcalc --version
```

## Usage

## 📝 Scripts & Automation

Below is an overview of the helper scripts included in this project, and how they integrate into the development and CI/CD processes:

- **`giti`**  
  A Bash script that automates the end-to-end release workflow:
  - Commits changes to `master` with a user-provided message.
  - Calculates and creates the next semantic version tag based on existing Git tags.
  - Pushes the new tag to GitHub.
  - Creates or updates the GitHub Release and uploads built binaries via `gh release`.
  - Records each created tag in `.giti_tag` for version tracking.

- **`build`**  
  A build helper script that compiles the `netcalc` binary for all supported platforms (Windows, Linux, macOS Intel, and macOS Apple Silicon) and places the executables in the `bin/` directory, ready for packaging or release.

- **`test_deploy.sh`**  
  Deploys the newly built binaries to a staging or test environment. This script can be used in CI to verify deployment steps before publishing a release.

- **`smoke_test.sh`**  
  Runs a suite of smoke tests against a deployed instance of `netcalc` (or its API/CLI), ensuring basic functionality works as expected immediately after deployment.

- **`setup_envrc.sh`**  
  Sets up environment variables and loads them via [`direnv`](https://direnv.net/) or similar tooling, making local development environment configuration automatic when entering the project directory.

Feel free to customize or extend these scripts to fit your workflow and environment.

## 🛠️ Useful Commands and Tools for This Project

Below are some Git commands and other tools that have helped resolve common issues encountered during development, build, and release workflows:

- `git fetch origin --tags --force`  
  *Fixes conflicts with existing tags on the remote when pushing new tags, ensuring the CI/CD pipeline picks up the correct tag for release.*

- `git reset --hard origin/main`  
  *Resets the local `main` branch to match the remote `origin/main`, discarding divergent local commits and resolving
  push rejections.*

- `git push origin main --force-with-lease`  
  *Safely force-pushes the local `main` branch to the remote, only if no one else has updated upstream, preventing
  “non-fast-forward” errors.*

- `git pull --rebase origin main`  
  *Rebases local commits on top of the latest remote `main`, keeping history linear and avoiding merge conflicts.*

- `go mod tidy`  
  *Generates missing `go.sum` entries and removes unused dependencies, fixing build errors related to missing modules.*

- `go test ./... -coverprofile=coverage.out` and Codecov  
  *Produces a coverage report and uploads it to Codecov, ensuring test coverage remains above the configured threshold.*

- `golangci-lint run --config .golangci.yml`  
  *Runs static analysis and code style checks, catching unused variables, formatting issues, and other lint warnings.*

- `gh release upload <tag> bin/*/* --clobber`  
  *Attaches all compiled binaries to a GitHub Release in one command, automating asset uploads.*


### 🔖 Tags and Releases Management

- **Create a new tag manually**
  ```bash
  git tag v1.0.0
  git push origin v1.0.0
  ```

- **Create a release from a tag using GitHub CLI**
  ```bash
  gh release create v1.0.0 --title "Release v1.0.0" --notes "Release description or changelog"
  ```

- **Remove a local and remote tag**
  ```bash
  git tag -d v0.0.1-test
  git push --delete origin v0.0.1-test
  ```

- **Remove a release from GitHub using GitHub CLI**
  ```bash
  gh release delete v0.0.1-test --yes
  ```

These commands are helpful when managing incorrect tags, cleaning up testing releases, or fixing mistakes in published versions.
