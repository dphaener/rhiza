# run-status-live-poll-waits-for-hold-verified

evidence · `01M7RJ07BGF2EC11K6PYMW28WA` · actor [claude-code-3270ca18](../actor/claude-code-3270ca18.md)

## result

PASS

## statement

Verified at 4c66d4717 in a detached worktree; every go test in systemd-run --user --scope -p MemoryMax=6G -p TasksMax=300, -p 2, RHIZA_*/SHIP_*/REX_* unset, one at a time. git show: the only code change is the poll line (run-status.txtar:35), no Go files, bound 400 and interval 0.025 unchanged, no assertion touched. (1) go test -trimpath -p 2 -count=3 ./cmd/rhiza -run 'TestScripts/^run-status$': exit 0. (2) scratch step first run = "sleep 0.5" with the script at abf1d4cc3: exit 1, FAIL at run-status.txtar:38 with '- Current step: first' in stdout; with the fix under the same delay, -count=3: exit 0. (3) same delay, poll grepping only '^- Current step: hold$': exit 1, FAIL at run-status.txtar:37 (status read the earlier stopped run). Scratch edits reverted, worktree removed. (4) go test -trimpath -p 2 -v ./cmd/rhiza -run 'TestScripts/run-status': exit 0, PASS for run-status, -interventions, -messages, -stalled, -tokens.

## Links

### Out

- verifies [run-status-live-poll-waits-for-hold](../capability/run-status-live-poll-waits-for-hold.md)

## History

- `01M7RHC6VMMBXVWS3BKKFF6F8N` Fix flaky test TestScripts/run-status: failed then passed on a rerun in PR #657 at head f5b2f2b99b740b4a9a18311f074ec9e4e2df94f2; reproduce it locally, repeatedly, then fix it or quarantine it
