# shipped-skip-reason

capability · `01M69DHDGP588ZJPVVXZ80QD8Z` · parent [daemon-skip-reasons](../daemon-skip-reasons.md)

## acceptance

A shipped case in the daemon's Skips table test and a step in daemon-status.txtar with a hand-written succeeded ship-change journal and fake rex show reason shipped; describe-golden.txtar is regenerated

## statement

rhiza daemon status reports such a change set with skip reason shipped

## Links

### Out

- decided-by [daemon-skips-shipped-rule](../../../../../../decision/daemon-skips-shipped-rule.md)

### In

- verifies [shipped-skip-reason-verified](../../../../../../evidence/shipped-skip-reason-verified.md)
- verifies [shipped-skip-reason-qa-1](../../../../../../evidence/shipped-skip-reason-qa-1.md)
- verifies [pr-139-daemon-skips-just-shipped](../../../../../../evidence/pr-139-daemon-skips-just-shipped.md)

## History

- `01M69DFHXFG7197RGNYYPWS4A5` The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened
