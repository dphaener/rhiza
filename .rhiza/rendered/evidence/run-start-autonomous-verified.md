# run-start-autonomous-verified

evidence · `01M5XCZZ3V8W2PY8N0AVPFB5CZ` · actor [claude-code-54c93d5c-verifier](../actor/claude-code-54c93d5c-verifier.md)

## result

PASS on 2026-10-26 at commit 9b213db (code as of 8b72ad5) on branch cs/autonomous-mode-and-the, re-run with -count=1 by an independent verifier (claude-code-54c93d5c-verifier)

## statement

Testscript run-autonomous starts a two-agent workflow with --autonomous on fake rex and claude: a dry run prints "autonomous": true and launches nothing; the journal's started entry has "autonomous":true; run status --json gives "autonomous": true and the markdown '- Autonomous: yes'; the fake agent logs RHIZA_AUTONOMOUS=1; after the wait stop, run resume without the flag launches step two with RHIZA_AUTONOMOUS=1. A run started without the flag (even with RHIZA_AUTONOMOUS=1 in the terminal) has no autonomous in the journal, status or agent env. describe-golden holds the new autonomous request and response fields and passes. TestFoldAutonomous folds the flag from the started entry. Hands-on: the worktree build, in a scratch repo with RHIZA_* unset, after rhiza init, change new 'X' and commits, run start --autonomous --permission-mode manual --dry-run --json printed "autonomous": true, exit 0.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/run-autonomous$'; go test -count=1 ./internal/run -run TestFoldAutonomous; go test -count=1 ./cmd/rhiza -run TestScripts/describe-golden (in make check); hands-on: go build -o $TMPDIR/rhiza-ver ./cmd/rhiza; rhiza run start --autonomous --permission-mode manual --dry-run --json

## Links

### Out

- verifies [run-start-autonomous](../capability/rhiza/workflows/autonomous-runs/run-start-autonomous.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
