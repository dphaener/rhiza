# daemon-state-at-start

capability · `01M6CAACNGV26NMDMKTZZ3759E` · parent [daemon-reexec-new-binary](../daemon-reexec-new-binary.md)

## acceptance

A Go test seeds a state file from an earlier image of the same pid with a last poll, max runs 3, one live run and one skip, takes the lock as the new image, writes the start state and runs daemon status before any poll: it reports the new start time, the old last poll, the request's max runs, the seeded run and skip, and no stopped line

## statement

A daemon image writes its state before its first poll, carrying the previous image's last poll, live runs, skips and unjournaled holds, so daemon status shows a running daemon's state through a re-exec

## Links

### Out

- decided-by [daemon-writes-state-at-start](../../../../../decision/daemon-writes-state-at-start.md)

### In

- verifies [daemon-state-at-start-verified](../../../../../evidence/daemon-state-at-start-verified.md)
- verifies [daemon-state-at-start-qa-1](../../../../../evidence/daemon-state-at-start-qa-1.md)
- verifies [pr-186-status-tells-the-truth](../../../../../evidence/pr-186-status-tells-the-truth.md)

## History

- `01M6CA9QMHN6N95AS3F935ZVN2` daemon status keeps the daemon's state through a re-exec, and a terminal current change set comes with the command that clears it
