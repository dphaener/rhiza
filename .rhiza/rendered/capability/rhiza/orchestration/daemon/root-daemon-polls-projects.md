# root-daemon-polls-projects

capability · `01M6RN85A2G5E84BYYP5201P8A` · parent [daemon](../daemon.md)

## acceptance

Daemon tests with a root of two initialised projects and one bare checkout start a ready change set in each initialised project and record the bare one as skipped, not initialised.

## statement

Each poll of the root daemon reads every registered, initialised project and starts its ship-change runs, plan requests and resumes as the per-repo daemon did, in that project's own checkout and journal; a project that is not initialised is skipped and daemon status says so.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)
- decided-by [m29-unready-listed](../../../../decision/m29-unready-listed.md)

### In

- verifies [root-daemon-polls-projects-verified](../../../../evidence/root-daemon-polls-projects-verified.md)
- verifies [root-daemon-polls-projects-qa-1](../../../../evidence/root-daemon-polls-projects-qa-1.md)
- verifies [root-daemon-polls-projects-qa-2](../../../../evidence/root-daemon-polls-projects-qa-2.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
