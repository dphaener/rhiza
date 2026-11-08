# watchdog-plan-interview

source · `01M6WJ1GK03S6X3QRQEP6KGNTT`

## quotes

```json
[]
```

## raw_sha256

43a55ed74fb45e332b148137e8a274b338095d3f032d320063459a73af75ade4

## session

fa29a4a8-96f3-4ec1-b17f-01a1b91c9b4a

## summary

Autonomous plan run: no owner present. The planner read the intent, the daemon code (daemonPoll1, daemonNudge, session.StallAfter, inbox stalled items) and the overlapping queued bug requests 01M6W8KN, 01M6W9VN and 01M6WBVX, and raised four choices as one question for the arbiter: scope against those requests, stalled-step recovery, PR recovery and shipped defaults. The arbiter answered with Decision watchdog-scope-and-defaults: the watchdog is the daemon-level backstop that acts only on runs that stalled or ended on a transient signature, never while a step's own recovery is live, with its own per-run retry budget; stalled autonomous steps are nudged, then stopped and retried, while owner-attended runs get only a log line and an inbox item; a PR with no CI gets origin/main merged in and pushed normally (never rebase or force-push), a green unmerged PR gets its run resumed so the workflow merges it. The defaults are machine keys rhiza.daemon.watchdog.*, written as flat full-key strings under [daemon] in the root's rhiza.toml with camelCase names: agentNudge 20m then agentRetry 40m for agent steps, commandLimit 30m for command steps, waitLimit 90m for CI and PR waits, retryBudget 2, ciReruns 1 (once per PR head), pollBehind 3 (warn past 3x the interval), fileBugs true, autonomousOnly true, enabled true. Project keys rhiza.project.watchdog.* of the same names (all but pollBehind) in the project's git config override the machine keys. The change set serves milestone-25 (the daemon builds and ships the queued requests on its own), as its sibling requests do.

## title

Daemon watchdog plan, autonomous interview with the arbiter

## Links

### In

- derived-from [watchdog-scope-and-defaults](../decision/watchdog-scope-and-defaults.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
