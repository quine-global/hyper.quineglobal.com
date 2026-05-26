#!/bin/sh
# Force an immediate re-fetch of GitHub release assets.
# Usage (from host):   ./scripts/refresh-releases.sh [http://host:8081]
# Usage (in container): ./scripts/refresh-releases.sh
set -e

url="${1:-http://localhost:8081}/internal/refresh"

if command -v curl >/dev/null 2>&1; then
    curl -fsS -X POST "$url"
elif command -v wget >/dev/null 2>&1; then
    wget -q -O- --post-data='' "$url"
else
    echo "error: curl or wget required" >&2
    exit 1
fi

echo
echo "Releases refreshed."
