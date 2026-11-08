# pr-375-daemon-watchdog

evidence · `01M6X1EZM2MWFV9MANZTWV6NFE` · actor [dphaener](../actor/dphaener.md)

## head

3cd60b3141138ad578bb68d5e5bfa4263f916397

## pr

https://github.com/dphaener/rhiza/pull/375

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #375 (https://github.com/dphaener/rhiza/pull/375) passed CI at 3cd60b3

## status

verified

## Links

### Out

- verifies [daemon-watchdog](../capability/rhiza/orchestration/daemon/daemon-watchdog.md)
- verifies [watchdog-policy](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-policy.md)
- verifies [watchdog-stalled-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stalled-step.md)
- verifies [watchdog-transient-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-transient-failure.md)
- verifies [watchdog-ci-rerun](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-ci-rerun.md)
- verifies [watchdog-pr-no-ci](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-pr-no-ci.md)
- verifies [watchdog-green-pr-unmerged](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-green-pr-unmerged.md)
- verifies [watchdog-stuck-requests](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stuck-requests.md)
- verifies [watchdog-retry-budget](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-retry-budget.md)
- verifies [watchdog-files-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-files-bugs.md)
- verifies [watchdog-poll-behind](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-poll-behind.md)
- verifies [watchdog-actions-shown](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-actions-shown.md)
- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
