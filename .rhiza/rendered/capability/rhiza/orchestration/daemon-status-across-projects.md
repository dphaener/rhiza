# daemon-status-across-projects

capability · `01M6RN85TXNPGZEANC6R3HFM27` · parent [orchestration](../orchestration.md)

## acceptance

A testscript with a root of two projects gets from rhiza daemon status --json at the root the root daemon once with each project's live runs, last starts and skip reasons, not initialised included, and only one project's part with --project and from inside that project.

## statement

rhiza daemon status run at the root (or with --all) covers every registered project, each item tagged with its project; run inside a project or with --project <name> it covers only that project.

## Links

### Out

- decided-by [m29-cli-and-tui-now](../../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-one-root-daemon](../../../decision/m29-one-root-daemon.md)

### In

- verifies [daemon-status-across-projects-verified](../../../evidence/daemon-status-across-projects-verified.md)
- verifies [daemon-status-across-projects-qa-1](../../../evidence/daemon-status-across-projects-qa-1.md)
- verifies [daemon-status-across-projects-verified-2](../../../evidence/daemon-status-across-projects-verified-2.md)
- verifies [daemon-status-across-projects-qa-2](../../../evidence/daemon-status-across-projects-qa-2.md)
- verifies [pr-342-cross-project-cli](../../../evidence/pr-342-cross-project-cli.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
