#!/bin/sh
# Updates the bundled C files from the Go module cache.
# Usage: go generate ./server/vectorstore/sqlitevec/
set -e

SQLITE_VEC_VERSION=$(go list -m -json github.com/asg017/sqlite-vec-go-bindings | grep '"Version"' | sed 's/.*"Version": "\(.*\)".*/\1/')
SQLITE_VEC_DIR=$(go list -m -json github.com/asg017/sqlite-vec-go-bindings | grep '"Dir"' | sed 's/.*"Dir": "\(.*\)".*/\1/')/cgo

MATTN_DIR=$(go list -m -json github.com/mattn/go-sqlite3 | grep '"Dir"' | sed 's/.*"Dir": "\(.*\)".*/\1/')

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Copying sqlite-vec ${SQLITE_VEC_VERSION} C source..."
cp "${SQLITE_VEC_DIR}/sqlite-vec.c" "${SCRIPT_DIR}/sqlite-vec.c"
cp "${SQLITE_VEC_DIR}/sqlite-vec.h" "${SCRIPT_DIR}/sqlite-vec.h"

echo "Copying sqlite3.h from mattn/go-sqlite3..."
cp "${MATTN_DIR}/sqlite3-binding.h" "${SCRIPT_DIR}/sqlite3.h"

echo "Done. Remember to update the version comment in vec.go."
