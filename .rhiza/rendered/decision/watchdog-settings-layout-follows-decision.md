# watchdog-settings-layout-follows-decision

decision · `01M6WK7YZZBT82TKQ4P9TJXNCG`

## reasoning

The root's rhiza.toml today holds daemon settings as flat full-key strings under [daemon] ("rhiza.daemon.maxRuns" = "3", "rhiza.daemon.selfUpdate", ...), settingDefs declares them that way, and daemonapp pairs rhiza.daemon.maxRuns with rhiza.project.maxRuns in git config. The Decision and capability follow that existing convention and are what the leaves will build and test against (TestWatchdogPolicy names [daemon] and rhiza.project.watchdog.*). A [rhiza.daemon.watchdog] table with snake_case file_bugs would be a new, inconsistent layout no code reads. The Source's summary is the planner's paraphrase of the arbiter's answer and lost detail; it is a provenance record, so it must describe the Decision it backs accurately rather than the Decision bending to it.

## statement

The Decision @watchdog-scope-and-defaults and the capability @watchdog-policy are authoritative on where the watchdog policy lives and what its keys are called: machine keys rhiza.daemon.watchdog.* as flat full-key strings under [daemon] in the root's rhiza.toml, camelCase names (fileBugs, agentNudge, retryBudget, ...), and project keys rhiza.project.watchdog.* of the same names (all but pollBehind) in the project's git config, overriding the machine keys. The Source @watchdog-plan-interview is corrected to match: its summary is rewritten to state that layout, the camelCase key names and the project keys, in place of '[rhiza.daemon.watchdog]' and 'file_bugs'. The Decision and capability are not changed.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6whrd-r3-f2](../question/review-01m6whrd-r3-f2.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
