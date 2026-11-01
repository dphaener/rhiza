# wait-stops-on-75

capability · `01M6A0NEQCRN4YPDW4EBSKYSKJ` · parent [wait-step](../wait-step.md)

## acceptance

A Go test in internal/run: a wait whose condition reports 75 with a stop: line ends the run stopped with that detail and the step not ended, and a resume polls the wait again; an exit other than 0 or 75 still reads as not held

## statement

A wait step whose condition exits 75 stops its run for the owner instead of polling on: the step is not ended, the run's stop detail is the last 'stop: ' line of the condition's output, and resuming the run runs the wait again

## Links

### Out

- decided-by [wait-settles-conflict-rule](../../../../../decision/wait-settles-conflict-rule.md)

### In

- verifies [wait-stops-on-75-verified](../../../../../evidence/wait-stops-on-75-verified.md)
- verifies [wait-stops-on-75-qa-1](../../../../../evidence/wait-stops-on-75-qa-1.md)
- verifies [pr-153-wait-settles-conflicting-pr](../../../../../evidence/pr-153-wait-settles-conflicting-pr.md)

## History

- `01M6A0JJDJ0FXD58K4WV88FV0A` A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict
