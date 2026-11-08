# watchdog-stalled-step-qa-1

evidence · `01M6WZESBF1N0PS24MA2PRBD1P` · actor [claude-code-9eef2845](../actor/claude-code-9eef2845.md)

## result

FAIL

## statement

A workflow command step 'sleep 300' with stall = "5s", run --autonomous, then rhiza daemon --once after 9s. The daemon logged 'step one silent since ...: stopped and resumed with --retry (budget 1 of 2)', the run showed stopped, and the next poll resumed it. An attended run of the same workflow was only logged ('attended, left to the owner') and listed in the inbox, with no stop. But stopping the run SIGTERMs only the run process: the step's child 'sh -c exec 2>&1 sleep 300' was reparented to PID 1 and kept running, and the retry started a second copy beside it.

## Links

### Out

- verifies [watchdog-stalled-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stalled-step.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
