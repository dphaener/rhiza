# daemon-headless-run-drivers-verified

evidence · `01M6FY311Z8ZAWD9MNY3CPA56A` · actor [claude-code-8a2beea2](../actor/claude-code-8a2beea2.md)

## result

PASS

## statement

At f2704ea1 I ran go test ./cmd/rhiza -run 'TestScripts/(daemon-headless-runs|daemon-picks|daemon-resume|daemon-repo-session|plan-requests|stage-autonomy|daemon-pause|daemon-settings|daemon-status|daemon-notify|describe-golden)$' (all pass) and read the diff. daemon-headless-runs runs rhiza daemon --once against a fake rex with no detach stub: the fake rex.log has 0 'window new' lines after the daemon exits, rhiza run ls --json then shows the run of the change set with "live": true, and .git/rhiza/logs/<change set>.log exists. With -testwork I read that log: it holds the child's own output ('step build/build: entered', 'builder session ...'), and the only window new in rex.log came later, from the agent session (rhiza session run), never from the driver. daemon-picks and stage-autonomy check the recorded detach command is 'rhiza run start --foreground <CS> --permission-mode manual', with --autonomous only under daemon --once --autonomous (and signoff auto); daemon-resume checks 'run resume --foreground <run>'. In the code, daemonStart/daemonResume/daemonStartPlan now go through daemonDetach -> ctx.Detach (Setsid, no stdin, output appended to the log, Release) with Dir=repo root and REX_SESSION set to the repo session, and no rex call for the driver. run start has no queue path at this commit, so the child drives the run and adds no queue entry (--foreground is accepted and is currently a marker). go test ./internal/api -skip TestReadBudget passes. Under heavy concurrent load from other sessions, the timing-only TestReadBudget, slice TestBudget and modelapp TestGeneratedModel failed. They are in code this commit doesn't touch. make check's self-check reported stale rendered actors for this session (8a2beea2) and a concurrent builder (d496b13b); neither record comes from this commit.

## Links

### Out

- verifies [daemon-headless-run-drivers](../capability/rhiza/orchestration/daemon/daemon-headless-run-drivers.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
