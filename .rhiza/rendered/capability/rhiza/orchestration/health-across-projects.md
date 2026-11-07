# health-across-projects

capability · `01M6RN85RW7YTX1MJ53SF5RB1J` · parent [orchestration](../orchestration.md)

## acceptance

A testscript with a root of two projects gets from rhiza health --json at the root one block per project plus the machine parts (daemon, disk, usage hold, machine trust) once, exiting 1 when any project's part is bad, and only one project's part with --project and from inside that project.

## statement

rhiza health run at the root (or with --all) covers every registered project, each item tagged with its project; run inside a project or with --project <name> it covers only that project.

## Links

### Out

- decided-by [m29-cli-and-tui-now](../../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-one-root-daemon](../../../decision/m29-one-root-daemon.md)

### In

- verifies [health-across-projects-verified](../../../evidence/health-across-projects-verified.md)
- verifies [health-across-projects-qa-1](../../../evidence/health-across-projects-qa-1.md)
- verifies [health-across-projects-verified-2](../../../evidence/health-across-projects-verified-2.md)
- verifies [health-across-projects-qa-2](../../../evidence/health-across-projects-qa-2.md)
- verifies [pr-342-cross-project-cli](../../../evidence/pr-342-cross-project-cli.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
