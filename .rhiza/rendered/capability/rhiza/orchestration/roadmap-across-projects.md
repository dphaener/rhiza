# roadmap-across-projects

capability · `01M6RN85QCT4TJYR8BAR6KZ499` · parent [orchestration](../orchestration.md)

## acceptance

A testscript with a root of two projects gets from rhiza roadmap --json at the root each project's milestone groups and backlog under its project name, and only one project's part with --project and from inside that project.

## statement

rhiza roadmap run at the root (or with --all) covers every registered project, each item tagged with its project; run inside a project or with --project <name> it covers only that project.

## Links

### Out

- decided-by [m29-cli-and-tui-now](../../../decision/m29-cli-and-tui-now.md)
- decided-by [m29-one-root-daemon](../../../decision/m29-one-root-daemon.md)

### In

- verifies [roadmap-across-projects-verified](../../../evidence/roadmap-across-projects-verified.md)
- verifies [roadmap-across-projects-qa-1](../../../evidence/roadmap-across-projects-qa-1.md)
- verifies [roadmap-across-projects-qa-2](../../../evidence/roadmap-across-projects-qa-2.md)
- verifies [pr-342-cross-project-cli](../../../evidence/pr-342-cross-project-cli.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
