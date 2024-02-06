#!/bin/sh
#/ .github/scripts/test.sh runs test on each go module in go-veikkaus.
#/ Arguments are passed to each go test invocation.
#/
#/ When UPDATE_GOLDEN is set, all directories named "golden" are removed before running tests

set -ex

CDPATH="" cd -- "$(dirname -- "$0")/.."

echo "Running the tests"

if [ "$#" = "0" ]; then
    set -- -race -covermode atomic ./...
fi

if [ -n "$UPDATE_GOLDEN" ]; then
    find . -name golden -type d -exec rm -rf {} +
fi

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname)"

for dir in $MOD_DIRS; do
    # In future skip example dir
    echo "testing $dir"
    (
        cd "$dir"
        go test "$@"
    ) || FAILED=1
done

if [ -n "$FAILED" ]; then
    exit 1
fi