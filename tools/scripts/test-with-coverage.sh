# Copyright Microsoft Corporation 2026
# SPDX-License-Identifier: MPL-2.0

# Wrapper for go test that writes coverage to a unique temp file per run.
#
# When gotestsum reruns failed tests, it re-invokes go test with the same flags.
# If -coverprofile is among them, the rerun (which only covers a subset of
# packages) overwrites the full coverage file, causing artificially low numbers.
#
# This script writes each run's coverage to a unique file (coverage.run.XXXXX.out).
# After gotestsum finishes, use gocovmerge to merge them into a single coverage.out.
#
# Usage with gotestsum:
#   gotestsum --raw-command --rerun-fails=5 --packages ./... -- \
#     bash ./tools/scripts/test-with-coverage.sh -run "^TestFoo" -timeout 90m
#   gocovmerge coverage.run.*.out > coverage.out

set -eu

COVFILE="coverage.run.$$.out"
exec go test -json -coverprofile="$COVFILE" -covermode=atomic ${COVERPKG:+-coverpkg="$COVERPKG"} "$@"
