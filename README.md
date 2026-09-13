# 📡 portwatch

[![CI](https://github.com/efelix32/portwatch/actions/workflows/ci.yml/badge.svg)](https://github.com/efelix32/portwatch/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/efelix32/portwatch)](https://goreportcard.com/report/github.com/efelix32/portwatch)
![License](https://img.shields.io/badge/license-MIT-green)

> A live, terminal-based view of every TCP/UDP connection on your machine — local/remote address, state, PID and process name, refreshed every 2 seconds. No `netstat` flags to memorize.

## Why

`netstat` and `ss` are powerful but their output is a wall of text you have to re-run and re-read. `portwatch` keeps a sorted, auto-refreshing table on screen instead, so you can watch a connection change state (`SYN_SENT` → `ESTABLISHED` → `TIME_WAIT`) in real time while debugging a flaky client or checking what a process is actually talking to.

## Install

```bash
go install github.com/efelix32/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/efelix32/portwatch.git
cd portwatch
go build -o portwatch .
```

## Usage

```bash
portwatch
```

```
portwatch -- live local connections

 PROTO  LOCAL                    REMOTE                   STATUS       PID    PROCESS
 tcp    0.0.0.0:5432             0.0.0.0:0                LISTEN       1204   postgres
 tcp    192.168.1.42:51231       142.250.72.10:443        ESTABLISHED  8891   chrome
 udp    0.0.0.0:5353             0.0.0.0:0                             612    mdnsd

q to quit -- refreshes every 2s
```

Cross-platform (Linux, macOS, Windows) via [gopsutil](https://github.com/shirou/gopsutil); the table UI is built with [bubbletea](https://github.com/charmbracelet/bubbletea).

## How it works

```
internal/scanner  -- reads OS connection table via gopsutil, sorts by port
internal/ui       -- bubbletea model: ticks every 2s, re-scans, redraws the table
main.go           -- wires the two together
```

`scanner` has no dependency on the terminal UI and is unit tested directly with fake connection data — no real network calls happen in tests. `ui`'s row-formatting logic is tested the same way.

## Development

```bash
go build ./...
go vet ./...
go test ./... -v
```

## License

MIT — see [LICENSE](LICENSE).
