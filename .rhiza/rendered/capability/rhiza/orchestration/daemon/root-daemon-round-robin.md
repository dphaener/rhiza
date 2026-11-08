# root-daemon-round-robin

capability · `01M6RN85B7BAZEZNH96W1P6D3A` · parent [daemon](../daemon.md)

## acceptance

internal/daemon pick tests: with maxRuns 3 and projects A (3 ready), B (2 ready) and C (1 ready), the picks are A1, B1, C1; with A's per-project max 1 and maxRuns 4, A gets one slot and the rest go to B and C.

## statement

The root daemon fills the machine-wide maxRuns round-robin across projects, each turn taking that project's next ready change set in its own roadmap order, except that every project's queued continuations take freed slots first (see root-continuations-first), and never runs more of a project's change sets at once than its optional per-project max.

## Links

### Out

- decided-by [m29-round-robin-slots](../../../../decision/m29-round-robin-slots.md)
- decided-by [continuations-first-machine-wide](../../../../decision/continuations-first-machine-wide.md)

### In

- verifies [root-daemon-round-robin-verified](../../../../evidence/root-daemon-round-robin-verified.md)
- verifies [root-daemon-round-robin-qa-1](../../../../evidence/root-daemon-round-robin-qa-1.md)
- verifies [root-daemon-round-robin-verified-2](../../../../evidence/root-daemon-round-robin-verified-2.md)
- verifies [root-daemon-round-robin-qa-2](../../../../evidence/root-daemon-round-robin-qa-2.md)
- verifies [root-daemon-round-robin-verified-3](../../../../evidence/root-daemon-round-robin-verified-3.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)
- verifies [pr-393-fix-a-queued-continuation-can](../../../../evidence/pr-393-fix-a-queued-continuation-can.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
