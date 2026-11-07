# m29-one-daemon

intent · `01M6RN37911APKX6WN4EB1CGR7` · parent [milestone-29](../milestone-29.md)

## statement

One root daemon builds every registered project's ready change sets, sharing run slots round-robin under one machine-wide cap and one usage governor, with machine settings in the root and project settings in each repo.

## Links

### Out

- decided-by [m29-one-root-daemon](../../decision/m29-one-root-daemon.md)
- decided-by [m29-round-robin-slots](../../decision/m29-round-robin-slots.md)
- decided-by [m29-settings-split](../../decision/m29-settings-split.md)
- decided-by [m29-root-required](../../decision/m29-root-required.md)
- decided-by [m29-unready-listed](../../decision/m29-unready-listed.md)
- decided-by [m29-rex-session-per-project](../../decision/m29-rex-session-per-project.md)

### In

- serves change set `01M6RN3RCGAZ9RNTWTZPKFM61M` (Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service)
- serves change set `01M6T1ZR3ZCDCYF96CSR2AYN7D` (Fix the usage governor's limits are wrong for the account: one limit hit (2026-11-04 10:45Z) taught both the session and the weekly limit, so the weekly limit became that week's running total (21.9M) and held every run for a day though the week reached 85.8M with no weekly cap hit; learned limits never relearn after a plan upgrade (the owner moved to Max 20x); and usageHoldOf counts a live run's tokens twice (already in the window's used, and again as a full run mean). A hit must teach only the kind it hit, learned limits must move up when use passes them without a hit, and a live run's projection must count only its remaining expected tokens)
- serves change set `01M6T2DCTWTJ0GW5AQMA2FZJF8` (Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available)
- serves change set `01M6TBDRB2N5MTW4BW0K2NZCG6` (Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
