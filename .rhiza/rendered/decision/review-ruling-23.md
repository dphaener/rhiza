# review-ruling-23

decision · `01M602SYYRSMM1S9CAT3J4J7DV`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

rhiza run approve [run] and rhiza run request-changes <file> [run] take the run by ref or the latest, exit 2 unless it is stopped at a signoff step, and resume the engine in-process; each requested change names a check, which re-runs. The built-in review.toml sets max_visits = 3 on its checks and signoff steps, and request-changes exits 2 when another visit would pass that cap.

## wording

the agent's choice in its plan, approved by the orchestrator acting as arbiter at plan approval

## Links

### In

- decided-by [review-signoff](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-signoff.md)
- decided-by [review-workflow-file](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-workflow-file.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
