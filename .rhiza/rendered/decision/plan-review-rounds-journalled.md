# plan-review-rounds-journalled

decision · `01M5XG42J0YXYCWW173RWBQ2P4`

## statement

Plan review rounds are counted from the run's journal: each rhiza run findings opens the next round, and a clean reviewer records an empty findings file; rhiza run review-check judges the latest round. The plan workflow loops review, critic and check, going to the gate on a pass and to revise then review on a fail, with max_visits 3 on review, critic and check and 2 on revise.

## wording

the orchestrator's ruling as arbiter for the owner, in the plan-workflow planning session

## Links

### In

- decided-by [plan-review-rounds](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-rounds.md)
- decided-by [plan-workflow-file](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-file.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
