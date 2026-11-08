# daemon-resume-once-defers-to-watchdog

decision · `01M6WK5SZATMGRYN3SWWDEFAK9`

## reasoning

The change set's intent is that the daemon recovers failed and stalled runs by resume --retry with a per-run retry budget, so the watchdog leaves resuming failed and given-up runs is the point of the work and must not be narrowed. daemon-resume-once is the older, landed rule and its 'leaves alone every failed run' clause was written before any recovery policy existed; its acceptance (one resume, a second interrupt gets none) still holds for the daemon's own plain resume. Amending daemon-resume-once to defer to the watchdog, as the finding and critic propose, removes the contradiction while keeping the inbox listing for runs the watchdog gives up on, which matches the intent's 'after the budget is spent leaves an inbox item' and the policy-with-shipped-defaults decisions (workflows-are-the-rules, warn-dont-safeguard).

## statement

This change set amends @daemon-resume-once so its one plain resume stays as it is, but failed runs and runs interrupted again after that resume are no longer left alone unconditionally: they are left to the watchdog's policy (@daemon-watchdog and its leaves), which may resume them with --retry within the per-run retry budget; only when the watchdog's policy does not apply, is turned off in rhiza.toml, or its budget is spent does the run stay listed in the inbox as failed or given up.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6whrd-r3-f1](../question/review-01m6whrd-r3-f1.md)

### In

- decided-by [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
