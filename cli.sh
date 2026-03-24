#!/usr/bin/env bash
# Chain Lens CLI entry point
# Usage:
#   ./cli.sh <fixture.json>              - Single transaction mode
#   ./cli.sh --block <blk.dat> <rev.dat> <xor.dat>  - Block mode

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Build the CLI binary if it doesn't exist or if source is newer
CLI_BIN="bin/chainlens-cli"

if [ ! -f "$CLI_BIN" ] || find cmd/cli internal -name "*.go" -newer "$CLI_BIN" 2>/dev/null | head -1 | grep -q .; then
    mkdir -p bin
    go build -o "$CLI_BIN" ./cmd/cli
fi

# Run the CLI with all arguments
exec "$CLI_BIN" "$@"
