name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.20

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: v1.64.8
          config: .golangci.yml

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.20

      - name: Run tests
        run: go test ./... -v
