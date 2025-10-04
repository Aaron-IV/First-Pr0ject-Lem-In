#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

usage() {
  echo "Usage: $0 [cli|serve] <file.txt> [addr]" >&2
  echo "  cli:   run CLI only (prints moves)" >&2
  echo "  serve: run server and open visualization link" >&2
  echo "Examples:" >&2
  echo "  $0 cli examples/example02.txt" >&2
  echo "  $0 serve examples/example02.txt :8080" >&2
}

if [[ ${1:-} == "" ]]; then
  usage; exit 1
fi

cmd="$1"; shift || true
file="${1:-}"; shift || true
addr="${1:-:8080}"

if [[ -z "$file" ]]; then
  usage; exit 1
fi

if [[ "$cmd" == "cli" ]]; then
  echo "Running CLI: $file" >&2
  go run ./cmd/main.go "$file"
elif [[ "$cmd" == "serve" ]]; then
  echo "Starting server: $file on $addr" >&2
  go run ./cmd/server "$file" "$addr"
else
  usage; exit 1
fi