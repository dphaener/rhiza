# run-start-through-daemon-verified-2

evidence · `01M6GHSDF0WVP9B7S959TE7N3R` · actor [claude-code-edcd0cbc](../actor/claude-code-edcd0cbc.md)

## result

PASS

## statement

Verifier for review cycle 2 at c19ff400, which fixes model finding f5. Read the diff: it adds only .rhiza log entries and rendered actor pages, with no Go code changed since 2d1f0db5. It amends roadmap-app-start-runs: p queues a plan run and s queues a ship run through run.start for the daemon, and the app says it queued the run. The amendment matches internal/tui/roadmapapp TestPlanRunIsQueuedForTheDaemon (one run.start with <cs> --workflow plan, view says 'Queued the plan run'), TestShipRunIsQueuedAutonomousWithAutoMode (<cs> --autonomous --permission-mode auto), TestNDeclinesAndCallsNothing and TestOFocusesLiveRunOrSaysNone. Re-ran the leaf's acceptance: go test ./internal/tui/roadmapapp -count=1 passes. go test ./cmd/rhiza -count=1 -run 'TestScripts/(run-start.*|daemon.*|describe-golden)$' passes, including run-start-queue (queues and says it waits for a daemon, daemon --once starts it headless within --max-runs 1, a backlog change set and a plan run on a draft are both started, --foreground drives in place) and daemon-headless-runs. make fmt-check vet self-check passes. In the full go test ./..., every package passes except the timing budgets internal/api TestReadBudget and internal/slice TestBudget, run at load average about 14 on 4 cores. TestBudget passes rerun alone. TestReadBudget times a generated model and stays over its 50ms medians under this load. Neither depends on this commit, which changes no Go code.

## Links

### Out

- verifies [run-start-through-daemon](../capability/rhiza/orchestration/daemon/run-start-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
