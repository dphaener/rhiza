# watchdog-scope-and-defaults

decision · `01M6WHZY78QBQEYRN9FEQWJ60A`

## reasoning

workflows-are-the-rules: recoveries that belong to a step stay in that step's workflow, and the watchdog only resumes runs, so the workflow's own steps (including merge) still decide. That makes it a backstop, not a second implementation of the three bug fixes, and it never merges outside a workflow. The intent says 'rebase', but bug 01M6WBVX and the repo's practice say merge origin/main, never rebase or force-push (force-push is denied here), so the watchdog merges. Owner-attended runs are left alone, as daemonNudgeStalled already does, because the owner may be the one the run waits on. The 20m nudge reuses today's daemonNudgeText threshold; doubling it before a retry gives the session one chance to answer. Budget 2 and one CI rerun match the intent's 'rerun failed CI once' and the bug requests' 'retry once', with one spare for a stall. warn-dont-safeguard: everything is policy with defaults in rhiza.toml, and it warns (inbox, log, health) rather than blocking.

## statement

(1) The watchdog is the daemon-level backstop; the three bug requests stay separate and ship on their own. Each adds a first-line recovery inside its workflow step (a step retries once on a vanished go $WORK, plan reruns failed CI once, ship-change merges main on a conflict). The watchdog acts only on what those miss: a run that has stalled, or has ended on a known transient signature. It never duplicates a step's recovery while the step is still live. Its resume counts against its own per-run retry budget, separate from any step's. (2) Stalled step (no journal entry or session report past its limit), autonomous runs only: first nudge the session (today's daemonNudge). If there is still no progress after a second limit, stop the step and run resume --retry, logging one line and charging the run's budget. Once the budget is spent, leave an inbox item naming what was tried. Owner-attended runs get only the log line and an inbox item, never a stop or retry. (3) A PR with no CI because it conflicts with main: the watchdog merges origin/main into the branch and pushes normally (never rebase or force-push), then records the push in the run's ship state. If the merge conflicts, it resumes the run so the run's own conflict handling (the 01M6WBVX fixer, once shipped) takes it; until then it leaves an inbox item. A green PR its run never merged: when the run ended failed at or before its merge step, the watchdog resumes it with --retry so the workflow's own merge step merges it, and never merges outside a workflow; a live run is left alone, and a run that ended succeeded (run resume refuses --retry on it) or whose merge needs the owner gets an inbox item instead. Every recovery caused by a bug files one deduplicated bug request with the evidence. (4) Shipped defaults, as settings in settingDefs like the other daemon settings, each key of one scope: machine keys rhiza.daemon.watchdog.* as flat full-key strings under [daemon] in the root's rhiza.toml, and project keys rhiza.project.watchdog.* of the same names (all but pollBehind) in the project's git config, as rhiza.daemon.maxRuns and rhiza.project.maxRuns are today. For a project the effective value is its project key when set, else the machine key, else the default; rhiza daemon settings shows each key's value, source and scope. A workflow step may set its own stall limit as a step field, which wins for that step and is shown with the workflow, not in the settings rows. Defaults: enabled = true; agentNudge = 20m (as today) and agentRetry = 40m for agent steps; commandLimit = 30m; waitLimit = 90m for CI and PR wait steps; retryBudget = 2 per run; ciReruns = 1 per PR head, counting a rerun the step already made; pollBehind = 3 (warn when a poll takes over 3x the poll interval); fileBugs = true; autonomousOnly = true.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-fa29a4a8-1](../question/question-fa29a4a8-1.md)
- derived-from [watchdog-plan-interview](../source/watchdog-plan-interview.md)

### In

- decided-by [daemon-watchdog](../capability/rhiza/orchestration/daemon/daemon-watchdog.md)
- decided-by [watchdog-policy](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-policy.md)
- decided-by [watchdog-stalled-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stalled-step.md)
- decided-by [watchdog-transient-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-transient-failure.md)
- decided-by [watchdog-ci-rerun](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-ci-rerun.md)
- decided-by [watchdog-pr-no-ci](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-no-ci.md)
- decided-by [watchdog-green-pr-unmerged](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-green-pr-unmerged.md)
- decided-by [watchdog-stuck-requests](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stuck-requests.md)
- decided-by [watchdog-retry-budget](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-retry-budget.md)
- decided-by [watchdog-files-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-files-bugs.md)
- decided-by [watchdog-poll-behind](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-poll-behind.md)
- decided-by [watchdog-actions-shown](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-actions-shown.md)
- decided-by [watchdog-pr-owned-by-run](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-owned-by-run.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
