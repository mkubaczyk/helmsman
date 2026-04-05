#!/usr/bin/env bash
# Runs three independent DSF files concurrently using the --parallel-files flag.
# Each DSF gets its own isolated temp dir and kubeconfig copy — no collisions.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Prefer the locally built binary, then fall back to whatever is on PATH.
if [ -z "${HELMSMAN_BIN:-}" ]; then
  if [ -x "$REPO_ROOT/helmsman" ]; then
    HELMSMAN="$REPO_ROOT/helmsman"
  else
    HELMSMAN="helmsman"
  fi
else
  HELMSMAN="$HELMSMAN_BIN"
fi

echo "Starting 3 parallel DSF executions via --parallel-files..."
echo ""

"$HELMSMAN" -parallel-files -p 3 -apply \
  -f "$SCRIPT_DIR/dsf-team-a.yaml" \
  -f "$SCRIPT_DIR/dsf-team-b.yaml" \
  -f "$SCRIPT_DIR/dsf-team-c.yaml"

EXIT_CODE=$?

echo ""
if [ $EXIT_CODE -eq 0 ]; then
  echo "All 3 DSF files applied successfully."
else
  echo "One or more DSF files failed (exit $EXIT_CODE)."
fi

exit $EXIT_CODE
