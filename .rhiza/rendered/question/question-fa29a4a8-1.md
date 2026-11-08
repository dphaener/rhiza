# question-fa29a4a8-1

question · `01M6WHXDC40GCZ2040HC7Y8W1S`

## session

fa29a4a8-96f3-4ec1-b17f-01a1b91c9b4a

## statement

Daemon watchdog scope: bug requests 01M6W8KN (retry a step once on a vanished go $WORK), 01M6W9VN (plan workflow reruns failed CI once) and 01M6WBVX (ship-change merges main on a rebase conflict) each add one recovery inside a workflow step. How should the watchdog relate to them? / When a run's step has had no journal entry or session report for longer than its limit, what does the watchdog do? / A run's PR has no CI (it conflicts with main) or is green but its run never merged it. What does the watchdog do? / Which shipped defaults should the watchdog's rhiza.toml policy have?

## Links

### In

- answers [watchdog-scope-and-defaults](../decision/watchdog-scope-and-defaults.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
