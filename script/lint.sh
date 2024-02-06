#!/bin/sh

set -e

GOLANGCI_LINT_VERSION="1.54.2"

CDPATH="" cd -- "$(dirname -- "$0")/.."
BIN="$(pwd -P)"/bin

mkdir -p "$BIN"

EXIT_CODE=0

fail() {
    echo "$@"
    EXIT_CODE=1
}

# Installing golangci-lint if bin/golangci-lint doesn't exist with the correct version
if ! "$BIN"/golangci-lint --version 2> /dev/null | grep -q "$GOLANGCI_LINT_VERSION"; then
    GOBIN="$BIN" go install "github.com/golangci/golangci-lint/cmd/golangci-lint@v$GOLANGCI_LINT_VERSION"
fi

MOD_DIRS="$(git ls-files '*go.mod' | xargs dirname |sort)"

for dir in $MOD_DIRS; do
    # In future skip example folder
    echo linting "$dir"
    (
        # Github Actions Output when running in an Action
        if [ -n "$GITHUB_ACTIONS" ]; then
            "$BIN"/golangci-lint run --path-prefix "$dir" --out-format github-actions
        else
            "$BIN"/golangci-lint run --path-prefix "$dir"
        fi
    ) || fail "failed linting $dir"
done

[ -z "$FAILED" ] || exit 1

exit "$EXIT_CODE"
