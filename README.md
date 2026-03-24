# Chain Lens

Chain Lens is a Bitcoin transaction and block analyzer written in Go.

It provides:
- A CLI for single transaction analysis and block file analysis.
- A web UI + API for interactive inspection.
- Structured JSON output designed for tooling and automation.

## Features

- Parse raw Bitcoin transactions from fixture JSON.
- Compute tx-level stats (`txid`, `wtxid`, fee, vbytes, script classification).
- Parse Bitcoin Core block/undo files (`blk*.dat`, `rev*.dat`) with XOR decoding key support.
- Export analysis reports to `out/`.
- Serve browser-based analysis views.

## Requirements

- Go 1.22+
- Linux/macOS shell environment

## Quick Start

```bash
./setup.sh
./cli.sh fixtures/tx/simple_p2wpkh.json
./web.sh
```

The web launcher prints a URL like `http://127.0.0.1:3000`.

## CLI Usage

Transaction mode:

```bash
./cli.sh <fixture.json>
```

Block mode:

```bash
./cli.sh --block <blk.dat> <rev.dat> <xor.dat>
```

Outputs are written to `out/`.

## Web API

- `GET /api/health`
- `POST /api/analyze`
- `POST /api/analyze/block`

## Project Layout

- `cmd/cli`: CLI entrypoint
- `cmd/web`: web server entrypoint
- `internal/`: parser/analyzer/core logic
- `fixtures/`: sample input data
- `web/static/`: frontend assets


