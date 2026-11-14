# default-tests-follow-conductor-qa-2

evidence · `01M7DTKQ89A6FZH9E725XYD5QR` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: run start --dry-run on a leafless accepted change set, with no --workflow and with --workflow ship-change, prints '- Workflow: ship-change' (exit 0). The four named scripts, and daemon-picks and daemon-status, pass. internal/daemon TestResumable|TestPick|TestSkips pass.

## Links

### Out

- verifies [default-tests-follow-conductor](../capability/rhiza/workflows/built-in-workflows/default-tests-follow-conductor.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
