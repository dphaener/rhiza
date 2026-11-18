# run-sessions-default-tmpdir-verified-2

evidence · `01M7R0M7BN400CXDSV19VXP3P2` · actor [claude-code-98f84b89](../actor/claude-code-98f84b89.md)

## result

PASS

## statement

Verified builder commit 5fff361b8 at branch head 55d323d37 in a clean detached worktree (no untracked files): make self-check exit 0 (the earlier stale rendered/index.md failure is gone); go test -trimpath -p 2 -count=1 ./internal/api/... -run 'SessionRun|Abandon|ShipRunEnv|Gotmp|RunTmp' exit 0 with TestSessionRunTmpdirDefault (no [env], and [env] without {run_tmp}: dry-run env TMPDIR=<run tmp>, folder absent on dry run, made on real start), TestSessionRunTmpdirOwn (TMPDIR=/x only, no run tmp made), TestAbandonRemovesSessionTmpdirFiles, TestShipRunEnvRunTmp passing; go test ./internal/workflow/... exit 0; rhiza ship generate changed 0 files; spec.go env Desc, docs/custom-workflows.md and docs/reference/workflow-file.md say TMPDIR defaults to {run_tmp} unless [env] sets it. Diff touches no internal/run code; no TODOs.

## Links

### Out

- verifies [run-sessions-default-tmpdir](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-sessions-default-tmpdir.md)

## History

- `01M7NWFRZTGPHG4DHNN8B6NWSP` When a run ends, whatever its outcome, nothing its steps wrote to the temporary folder is left behind: the daemon sets TMPDIR to the run's {run_tmp} for every process a run's steps start unless the workflow's [env] sets TMPDIR itself, so the agents and commands of a workflow whose [env] names no {run_tmp} write under {run_tmp}, which is already removed when the run ends; the /tmp/rhiza-projection-* leak stays bug 01M7M3JZFW's
