# daemon-version-in-health

capability · `01M6FESA0V3RCW576YB1TKM3VP` · parent [daemon](../daemon.md)

## acceptance

A testscript: daemon status --json carries the running daemon's version; health --json reports the daemon part warn naming both commits when origin/main is ahead of it, and ok when they match

## statement

rhiza daemon status shows the commit the running daemon was built from, and rhiza health warns on its daemon part when origin/main is ahead of that commit

## Links

### Out

- decided-by [m24-update-from-main](../../../../decision/m24-update-from-main.md)

### In

- verifies [daemon-version-in-health-verified](../../../../evidence/daemon-version-in-health-verified.md)
- verifies [daemon-version-in-health-qa-1](../../../../evidence/daemon-version-in-health-qa-1.md)
- verifies [pr-199-the-daemon-pulls-main-and-runs](../../../../evidence/pr-199-the-daemon-pulls-main-and-runs.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs
