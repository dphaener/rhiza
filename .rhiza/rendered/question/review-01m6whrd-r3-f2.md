# review-01m6whrd-r3-f2

question · `01M6WK4B4K1M2XTSZBBCN4F1C0`

## statement

Plan review round 3 left finding f2 (contradiction on watchdog-plan-interview): The Source's summary of the arbiter's answer puts the defaults 'under [rhiza.daemon.watchdog] in rhiza.toml' with 'file_bugs on', and mentions only machine settings. The Decision watchdog-scope-and-defaults it backs, and watchdog-policy, instead use flat full-key strings rhiza.daemon.watchdog.* under [daemon], a camelCase fileBugs key, and project keys rhiza.project.watchdog.* in git config that override the machine keys. The provenance record and the decision disagree on where the policy lives and what its keys are called. The critic confirmed it: The watchdog-plan-interview summary says the defaults live 'under [rhiza.daemon.watchdog] in rhiza.toml' with 'file_bugs on' and mentions no project keys. watchdog-policy states flat rhiza.daemon.watchdog.* keys under [daemon], a camelCase fileBugs key, and rhiza.project.watchdog.* project keys in git config. The Source behind the decision misstates its layout and key names. How should it be settled?

## Links

### In

- answers [watchdog-settings-layout-follows-decision](../decision/watchdog-settings-layout-follows-decision.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
