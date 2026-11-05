# question-58fb9cbb-1

question · `01M6MYBH597QRAM0PS1MTJFMEP`

## session

58fb9cbb-b516-4d6f-a418-01a1a9e44747

## statement

Change set 'Milestone planning always' serves no intent yet. Which landed intent should it serve? (Recent fix change sets like 01M6MT67 serve milestone-25, the open milestone; milestone-22 is the one about switching stages between autonomous and owner-driven.) / How should Rhiza make sure a milestone plan run (a record change set, or a request of kind milestone) never runs autonomously?

## Links

### In

- answers [milestone-plan-never-autonomous](../decision/milestone-plan-never-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
