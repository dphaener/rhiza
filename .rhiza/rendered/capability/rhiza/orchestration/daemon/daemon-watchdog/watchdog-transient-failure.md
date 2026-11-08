# watchdog-transient-failure

capability · `01M6WJ2CRTGK0ZMW94E2BA801H` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogTransient shows each signature in a failed step's output leading to one resume --retry, a failure matching no signature left for the owner, and a live step with the same output left alone

## statement

The daemon resumes with --retry a run that ended failed at a step whose output matches a known transient signature (a vanished go $WORK file, a held git index.lock, a network or GitHub 5xx error), within the run's watchdog budget, and never acts on a step that is still live

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- blocks [review-01m6whrd-r3-f1](../../../../../question/review-01m6whrd-r3-f1.md)
- verifies [watchdog-transient-failure-verified](../../../../../evidence/watchdog-transient-failure-verified.md)
- verifies [watchdog-transient-failure-qa-1](../../../../../evidence/watchdog-transient-failure-qa-1.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
