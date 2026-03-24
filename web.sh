#!/usr/bin/env bash
# Chain Lens Web Server entry point
# Starts the web visualizer on PORT (default: 3000)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Default port
PORT="${PORT:-3000}"
export PORT

# Build the web binary if it doesn't exist or if source is newer
WEB_BIN="bin/chainlens-web"

if [ ! -f "$WEB_BIN" ] || find cmd/web internal web -name "*.go" -newer "$WEB_BIN" 2>/dev/null | head -1 | grep -q .; then
    mkdir -p bin
    go build -o "$WEB_BIN" ./cmd/web
fi

# If the port is already in use, avoid a confusing bind error on repeated runs.
existing_pid="$(ss -ltnp "( sport = :$PORT )" 2>/dev/null | sed -n 's/.*pid=\([0-9]\+\).*/\1/p' | head -n1)"
if [ -n "$existing_pid" ]; then
    existing_cmd="$(ps -p "$existing_pid" -o comm= 2>/dev/null | tr -d '[:space:]')"
    if [ "$existing_cmd" = "chainlens-web" ]; then
        echo "http://127.0.0.1:$PORT"
        echo "chainlens-web is already running on port $PORT (pid $existing_pid)."
        exit 0
    fi

    echo "Port $PORT is already in use by '$existing_cmd' (pid $existing_pid)."
    echo "Use another port: PORT=3001 ./web.sh"
    exit 1
fi

# Run the web server (exec to replace shell process for signal handling)
exec "$WEB_BIN"
