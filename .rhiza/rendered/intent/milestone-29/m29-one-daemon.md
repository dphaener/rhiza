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

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
