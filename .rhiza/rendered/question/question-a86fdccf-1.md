# question-a86fdccf-1

question · `01M77N6XND5EV3P8BE2EBEDATK`

## session

a86fdccf-e6b8-4c91-ad56-01a1cf50f590

## statement

Change set 01M77N1 (Fix watchdog unmerged PR at land): ship-change's land step was refused because another accepted change set (auto-update) had landed a new statement for the same node first. land-triage only passes gate failures, so the run ended and the orchestrator repaired it by hand with `rhiza set --rebase` (commit 480b6acf3). How should ship-change handle a land refused for a model conflict?

## Links

### In

- answers [ship-land-model-conflict-rebase](../decision/ship-land-model-conflict-rebase.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
