# daemon-status-across-projects-verified

evidence · `01M6STBT2EC0M5JFEGX11NTNYR` · actor [claude-code-29a430a5](../actor/claude-code-29a430a5.md)

## result

PASS

## statement

Verifier at 1452562f. go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-status-projects|daemon-projects|describe-golden)$' passes. daemon-status-projects.txtar builds a root of two initialised projects (repo, second) plus a bare checkout, runs one root daemon poll that starts a run in each, then asserts: daemon status --json at the root shows the root daemon once (one running:false), a projects part for repo and for second each with its last_started, and bare once under skipped_projects with reason not initialised; the human text shows Project repo:/Project second: sections; --project repo and running from inside second give only that project's part, with no projects array; --all from inside a project covers both; an unknown --project is refused as not registered. I also rebuilt the scenario by hand with a scratch binary and a fake rex, adding a second accepted change set per project: at the root each project's part listed its own runs and skipped change sets with reasons (no-milestone) under its name, and bare was listed as not initialised; from inside second and with --project second only second's skips were shown. Diff read: runDaemonStatus reads each scope's daemon state into DaemonProjectStatus and empties the top-level runs/skips when merged. go test ./internal/api: four failures unrelated to this diff (only runDaemonStatus changed). TestReadBudget is a timing budget, with load average about 8. Three usage-limit tests (TestUsageWindowsHitKinds, TestLimitHitsBecomeLearnedLimits, TestSessionRunRecordsLimitHit) fail only inside this run checkout and pass on the same commit in a fresh worktree, so the cause is this checkout's environment.

## Links

### Out

- verifies [daemon-status-across-projects](../capability/rhiza/orchestration/daemon-status-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
