# plan-review-arbitrated

capability · `01M5XG434N2323BPMA11HAJ0JP` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript runs an autonomous workflow whose command step journals an asked entry for a question node it added: a fake arbiter answers it, the journal has arbiter and arbitrated entries for it, and the run goes on to its next step; with a failing arbiter the run fails naming the question; in a normal run no arbiter runs and the question stays open.

## statement

In an autonomous run, after a command step the engine reads the run's journal again and has the arbiter decide each asked question it has not yet decided, as it does for a session's question; if the arbiter fails, the run fails naming the question. A normal run leaves those questions open for the owner.

## Links

### Out

- decided-by [plan-review-left-arbitrated](../../../../../decision/plan-review-left-arbitrated.md)
- decided-by [arbiter-answers-where-question-lives](../../../../../decision/arbiter-answers-where-question-lives.md)

### In

- verifies [plan-review-arbitrated-verified](../../../../../evidence/plan-review-arbitrated-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
