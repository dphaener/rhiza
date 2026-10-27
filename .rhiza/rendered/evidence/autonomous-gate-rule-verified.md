# autonomous-gate-rule-verified

evidence · `01M5XCZZM99A0K8ENVQTS01QKG` · actor [claude-code-54c93d5c-verifier](../actor/claude-code-54c93d5c-verifier.md)

## result

PASS on 2026-10-26 at commit 9b213db (code as of 8b72ad5) on branch cs/autonomous-mode-and-the, re-run with -count=1 by an independent verifier (claude-code-54c93d5c-verifier)

## statement

Testscript autonomous-gate-rule proposes a change set and runs accept --rule autonomous: with no RHIZA_RUN it is usage, exit 2, and change show stays proposed; with RHIZA_RUN naming no run, and naming a non-autonomous run, it is usage, exit 2, state proposed and no accepted journal entry; with RHIZA_RUN naming the unfinished autonomous run it accepts, and change show --json gives rule autonomous and reason 'autonomous run <run id>', with one accepted entry (change, log_entry) in that run's journal. A finished autonomous run is refused again; reject is guarded the same way. TestFoldAccepts folds accepted entries. Hands-on: the worktree build in a scratch repo, RHIZA_RUN unset, accept --rule autonomous printed 'cannot accept under rule autonomous: ... RHIZA_RUN is not set' and exited 2.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/autonomous-gate-rule$'; go test -count=1 ./internal/run -run TestFoldAccepts; hands-on: rhiza accept --rule autonomous (RHIZA_RUN unset) -> exit 2

## Links

### Out

- verifies [autonomous-gate-rule](../capability/rhiza/workflows/autonomous-runs/autonomous-gate-rule.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
