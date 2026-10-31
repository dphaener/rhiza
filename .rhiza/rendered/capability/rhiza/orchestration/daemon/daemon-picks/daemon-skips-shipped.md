# daemon-skips-shipped

capability · `01M69DHC8T4MV7X7VWZFYD5YRP` · parent [daemon-picks](../daemon-picks.md)

## acceptance

Table cases in internal/daemon/pick_test.go: a succeeded ship-change run after the accept is skipped; a succeeded plan run is picked; a succeeded ship-change run followed by a reopen is picked

## statement

The daemon does not pick an accepted change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened

## Links

### Out

- decided-by [daemon-skips-shipped-rule](../../../../../decision/daemon-skips-shipped-rule.md)

### In

- verifies [daemon-skips-shipped-verified](../../../../../evidence/daemon-skips-shipped-verified.md)
- verifies [daemon-skips-shipped-qa-1](../../../../../evidence/daemon-skips-shipped-qa-1.md)
- verifies [pr-139-daemon-skips-just-shipped](../../../../../evidence/pr-139-daemon-skips-just-shipped.md)

## History

- `01M69DFHXFG7197RGNYYPWS4A5` The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened
