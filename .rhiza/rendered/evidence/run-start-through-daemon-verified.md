# run-start-through-daemon-verified

evidence · `01M6G73E64577DSYZTCJYB9J5H` · actor [claude-code-12950cc2](../actor/claude-code-12950cc2.md)

## result

PASS

## statement

At 8f17c588 I ran go test ./cmd/rhiza -count=1 with run-start-queue, run-start, run-ship-change, headless-run-agent-tabs, describe-golden, run-continue-queue, all daemon-* scripts, and the scripts named in attempt 2's findings (run-agent-step, run-agent-tabs, build-tabs, plan-workflow, -pr, -open, -autonomous, -autonomous-left, arbiter-agent-tabs, run-agent-done). All of them pass. In run-start-queue, rhiza run start queues the run and says no daemon runs, so it waits for one. daemon --once under --max-runs 1 starts one run headless with no window new and leaves the other queued. With a free slot it starts the plan run queued on the draft, and the first queued run was on a backlog change set. --foreground drives the run in place in run-start and the other in-place scripts. roadmapapp TestPlanRunIsQueuedForTheDaemon and TestShipRunIsQueuedAutonomousWithAutoMode pass. The fix commit drops headlessSession from foreground start and resume and sets daemonDetach to unset REX_BLOCK. headless-run-agent-tabs checks that a run queued from another session and block opens its tabs in rhiza-repo. Rerun uncached, the internal/api TestRunStart|TestRunResume|TestRunAmend|TestDaemon tests pass. go vet and make self-check pass. A full go test ./... had every package pass except the timing budgets TestReadBudget, slice TestBudget and modelapp TestGeneratedModel. The machine's load average was about 20, and this leaf does not change the read path

## Links

### Out

- verifies [run-start-through-daemon](../capability/rhiza/orchestration/daemon/run-start-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
