# inbox-across-projects

capability · `01M6RN85P5RQ03S1THHX2PE816` · parent [orchestration](../orchestration.md)

## acceptance

A testscript with a root of two projects gets from rhiza inbox --json at the root the merged list with a project field on every item, and inbox approve/open act on an item from the other project, and only one project's part with --project and from inside that project.

## statement

rhiza inbox run at the root (or with --all) covers every registered project, each item tagged with its project; run inside a project or with --project <name> it covers only that project.

## Links

### Out

- decided-by [m29-cli-and-tui-now](../../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-one-root-daemon](../../../decision/m29-one-root-daemon.md)

### In

- verifies [inbox-across-projects-verified](../../../evidence/inbox-across-projects-verified.md)
- verifies [inbox-across-projects-qa-1](../../../evidence/inbox-across-projects-qa-1.md)
- verifies [inbox-across-projects-verified-2](../../../evidence/inbox-across-projects-verified-2.md)
- verifies [inbox-across-projects-qa-2](../../../evidence/inbox-across-projects-qa-2.md)
- verifies [pr-342-cross-project-cli](../../../evidence/pr-342-cross-project-cli.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
