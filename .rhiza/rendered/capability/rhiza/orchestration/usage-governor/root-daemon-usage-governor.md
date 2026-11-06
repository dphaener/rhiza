# root-daemon-usage-governor

capability · `01M6RN85CZGYA2FJCNC52YWNAZ` · parent [usage-governor](../usage-governor.md)

## acceptance

internal/daemon usage tests: a limit hit recorded by one project's driver holds starts in a second project, and the hold lifts for both at the reset time.

## statement

The root daemon keeps one usage governor for the machine: run drivers of every project record usage-limit hits in the root's state folder, and the governor's hold stops new runs in all projects at once.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)

### In

- verifies [root-daemon-usage-governor-verified](../../../../evidence/root-daemon-usage-governor-verified.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
