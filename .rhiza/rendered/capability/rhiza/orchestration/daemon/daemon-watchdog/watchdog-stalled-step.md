# watchdog-stalled-step

capability · `01M6WJ2CKKVZMZ0EMK7HYXK1WH` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogStalled, with a fake clock, fake sessions and a fake run driver, shows an autonomous run nudged at the first limit, stopped and retried at the second with one log line and one budget charge, a run whose step makes progress in between left alone, and an owner-attended run getting a log line and an inbox item but no stop or retry

## statement

When a live autonomous run's step has had no journal entry or session report past its step kind's limit, the daemon nudges the session, and if nothing moves by the second limit it stops the step and resumes the run with --retry, charging the run's watchdog budget; an owner-attended run gets only a log line and an inbox item, never a stop or retry

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- verifies [watchdog-stalled-step-verified](../../../../../evidence/watchdog-stalled-step-verified.md)
- verifies [watchdog-stalled-step-qa-1](../../../../../evidence/watchdog-stalled-step-qa-1.md)
- verifies [watchdog-stalled-step-verified-2](../../../../../evidence/watchdog-stalled-step-verified-2.md)
- verifies [watchdog-stalled-step-qa-2](../../../../../evidence/watchdog-stalled-step-qa-2.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
