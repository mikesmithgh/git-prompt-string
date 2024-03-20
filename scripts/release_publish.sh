#!/usr/bin/env bash
set -euxo pipefail

release_notes="$1"
printf "%s" "$release_notes" >/tmp/release-notes.md
goreleaser release --clean --release-notes /tmp/release-notes.md
