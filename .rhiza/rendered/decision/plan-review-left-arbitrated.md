# plan-review-left-arbitrated

decision · `01M5XG42916B2D34M1PNYPVJPB`

## ruling

the arbiter decides each finding left after round 3, then the run accepts; if the arbiter can't, the run fails with a report

## statement

In an autonomous plan run, each finding left after the third review round is decided by the arbiter: rhiza run review-check journals each question it adds as an asked entry, the engine re-reads the journal after a command step and has the arbiter decide each asked question not yet decided, and the gate then accepts; if the arbiter can't decide one, the run fails naming the question and run status is the report.

## wording

the orchestrator's ruling as arbiter for the owner, in the plan-workflow planning session

## Links

### Out

- answers [plan-review-left-autonomous](../question/plan-review-left-autonomous.md)

### In

- decided-by [plan-review-arbitrated](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-arbitrated.md)
- decided-by [plan-review-rounds](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-rounds.md)
- decided-by [plan-gate](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
