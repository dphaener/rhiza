# tui-launch-usage

capability · `01M68587MGRTB20J74QQF0AZS6` · parent [rex-app-actions](../rex-app-actions.md)

## acceptance

A testscript runs ! exec rhiza tui launch and checks stderr has usage: rhiza tui app <name> and the exit code is 2

## statement

rhiza tui launch is no command: it exits 2 with the rhiza tui app <name> usage and opens no picker

## Links

### Out

- decided-by [m6-housekeeping-folded](../../../decision/m6-housekeeping-folded.md)

### In

- verifies [tui-launch-usage-verified](../../../evidence/tui-launch-usage-verified.md)
- verifies [tui-launch-usage-qa-1](../../../evidence/tui-launch-usage-qa-1.md)
- verifies [pr-122-remove-the-all-apps-picker](../../../evidence/pr-122-remove-the-all-apps-picker.md)

## History

- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.
