#!/usr/bin/env bash
# Decide whether CI must run the code jobs. Run in a checkout of the commit
# under test (HEAD) with full history.
#
#   ci-changes.sh BASE [GREEN]
#
# BASE is the pull request's base commit; GREEN is the pull request's latest
# earlier commit whose CI passed, or empty. Prints code=false when no path
# outside .rhiza/ differs between GREEN and HEAD, code=true otherwise. With no
# usable GREEN (empty, or not an ancestor of HEAD) the diff runs from the merge
# base with BASE instead.
set -euo pipefail

base=${1:?usage: ci-changes.sh BASE [GREEN]}
green=${2:-}

from=
if [ -n "$green" ] && git merge-base --is-ancestor "$green" HEAD 2>/dev/null; then
  from=$green
else
  from=$(git merge-base "$base" HEAD)
fi

outside=$(git diff --name-only "$from" HEAD | grep -v '^\.rhiza/' || true)
if [ -n "$outside" ]; then
  echo "code=true"
else
  echo "code=false"
fi
