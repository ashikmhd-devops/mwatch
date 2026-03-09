# mwatch — macOS System Monitor

[![Go](https://github.com/ashikmhd-devops/mwatch/actions/workflows/go.yml/badge.svg)](https://github.com/ashikmhd-devops/mwatch/actions/workflows/go.yml)

A fast, modern terminal resource monitor for macOS. Built with Go + Bubble Tea + Lip Gloss.

```
◆ mwatch  macmini.local  darwin/arm64  up 2 days, 3 hrs  ●···

┌─────────────────────────────────┐ ┌─────────────────────────────────┐
│ ▸ CPU                           │ │ ▸ Memory                        │
│  Overall ████████░░░░ 68.2%     │ │  RAM   ██████████░░ 11.2/16.0GB │
│  Load  1.92  1.45  1.20         │ │  Swap  ░░░░░░░░░░░░  0.1/2.0GB  │
│  P0  ██████░░░░  72.1%          │ │  Avail  4.8 GB                  │
│  P1  ████░░░░░░  48.3%  ...     │ └─────────────────────────────────┘
└─────────────────────────────────┘
┌─────────────────────────────────┐ ┌─────────────────────────────────┐
│ ▸ Network                       │ │ ▸ Disk  /                       │
│  ↑ Upload    12.3 KB/s          │ │  Used  ██████░░░░  280/512 GB   │
│  ↓ Download  341.0 KB/s         │ │  I/O   R: 0.12 MB/s | W: 0.03  │
└─────────────────────────────────┘ └─────────────────────────────────┘
┌──────────────────────────────────────────────────────────────────────┐
│ ▸ Processes                                                          │
│  PID    NAME                CPU%    MEM MB  STATUS  USER             │
│  ─────────────────────────────────────────────────────────          │
│  1234   node                12.3    234     S       ashik            │
│  ...                                                                 │
└──────────────────────────────────────────────────────────────────────┘
  c:sort CPU  │  m:sort MEM  │  p:sort PID  │  n:sort NAME  │  q:quit
```

## Requirements

- Go 1.22+
- macOS (Apple Silicon or Intel)

## Install

```bash
# Clone
git clone https://github.com/ashikmhd/mwatch
cd mwatch

# Fetch deps & build
make

# Install to /usr/local/bin
make install

# Run
mwatch
```

## Dev

```bash
# Run directly
make run

# Build Apple Silicon optimized binary
make build-arm64
```

## Key Bindings

| Key | Action |
|-----|--------|
| `c` | Sort processes by CPU |
| `m` | Sort processes by Memory |
| `p` | Sort processes by PID |
| `n` | Sort processes by Name |
| `q` | Quit |

## Design

- **Stack**: Go 1.22, Bubble Tea, Lip Gloss, gopsutil
- **Theme**: Industrial Precision — deep black `#0d0d0d`, electric cyan `#00d4ff`, amber warnings
- **Refresh**: Every 1.5 seconds (configurable)
- **Binary size**: ~8MB (stripped with `-ldflags="-s -w"`)
- **CPU usage of mwatch itself**: <0.5%
