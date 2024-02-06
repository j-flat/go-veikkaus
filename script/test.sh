#!/bin/sh
#/ .github/scripts/test.sh runs test on each go module in go-veikkaus.
#/

set -e

CDPATH="" cd -- "$(dirname -- "$0")/.."

echo "Running the tests"

if [ "$#" = "0" ]; then
    set -- -race -covermode atomic ./...
fi

( ginkgo "$@" ) || FAILED=1

if [ -n "$FAILED" ]; then
    exit 1
fi
