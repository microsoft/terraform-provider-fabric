# Copyright Microsoft Corporation 2026
# SPDX-License-Identifier: MPL-2.0

# Wrapper for `go test` that emits Go binary coverage data (covdata) into
# $GOCOVERDIR. When $GOCOVERDIR is set, `go test -cover` automatically writes
# covcounters.* / covmeta.* files there, so multiple test runs (including
# gotestsum reruns) accumulate coverage without overwriting each other.
#
# After tests finish, merge with:
#   go tool covdata textfmt -i="$GOCOVERDIR" -o=coverage.out
#
# Usage with gotestsum:
#   export GOCOVERDIR=.coverdata
#   gotestsum --raw-command --rerun-fails=5 --packages ./... -- \
#     bash ./tools/scripts/gotest-cover.sh -run "^TestFoo" -timeout 90m

set -eu

: "${GOCOVERDIR:?GOCOVERDIR must be set}"

# Resolve to an absolute path: each test binary runs in its package directory,
# so a relative path would not point to the shared coverage directory.
mkdir -p "$GOCOVERDIR"
ABS_GOCOVERDIR="$(cd "$GOCOVERDIR" && pwd)"

exec go test -json -cover -covermode=atomic ${COVERPKG:+-coverpkg="$COVERPKG"} "$@" -args -test.gocoverdir="$ABS_GOCOVERDIR"
