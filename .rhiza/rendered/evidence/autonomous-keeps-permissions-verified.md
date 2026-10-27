# autonomous-keeps-permissions-verified

evidence · `01M5XCZZC44X113SBPW84AFSP4` · actor [claude-code-54c93d5c-verifier](../actor/claude-code-54c93d5c-verifier.md)

## result

PASS on 2026-10-26 at commit 9b213db (code as of 8b72ad5) on branch cs/autonomous-mode-and-the, re-run with -count=1 by an independent verifier (claude-code-54c93d5c-verifier)

## statement

Testscript autonomous-permissions runs one workflow on a fake adapter without and with --autonomous; both launches log --permission-mode manual and neither has bypassPermissions or dangerously. run start --autonomous --permission-mode bypassPermissions exits non-zero with 'permission mode bypassPermissions bypasses the agent's permission checks: that needs a machine with trust = "dedicated"' (also in --json), and the runs directory still holds 2 journals and no third launch is logged.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/autonomous-permissions$'

## Links

### Out

- verifies [autonomous-keeps-permissions](../capability/rhiza/workflows/autonomous-runs/autonomous-keeps-permissions.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
