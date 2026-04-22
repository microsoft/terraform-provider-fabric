#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation
# SPDX-License-Identifier: MPL-2.0

# Wrapper for go test that prevents reruns from overwriting coverage data.
#
# When gotestsum reruns failed tests, it re-invokes go test with the same flags.
# If -coverprofile is among them, the rerun (which only covers a subset of
# packages) overwrites the full coverage file, causing artificially low numbers.
#
# This script collects coverage only on the first run. On subsequent runs
# (detected by coverage.out already existing), it skips -coverprofile so the
# original full-suite coverage data is preserved.
#
# Usage with gotestsum:
#   gotestsum --raw-command --rerun-fails=5 --packages ./... -- \
#     ./tools/scripts/test-with-coverage.sh -run "^TestFoo" -timeout 90m

set -eu

if [ ! -f "coverage.out" ]; then
  exec go test -json -coverprofile="coverage.out" -covermode=atomic ${COVERPKG:+-coverpkg="$COVERPKG"} "$@"
else
  exec go test -json "$@"
fi
