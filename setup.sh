#!/usr/bin/env bash
set -euo pipefail

# Install dependencies and prepare fixtures.

for gz in fixtures/blocks/*.dat.gz; do
  dat="${gz%.gz}"
  if [[ ! -f "$dat" ]]; then
    echo "Decompressing $(basename "$gz")..."
    gunzip -k "$gz"
  fi
done

go mod download

echo "Setup complete"
