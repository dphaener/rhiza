# review-ruling-18

decision · `01M602SWS47TH0HH7YY2WJA3H8`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

Fixers run one at a time because they share the run worktree: leaf fixers in leaf order, then the one change-set fixer, with a fixed leaf's verifier running beside the next fixer. Only the checks that failed re-run, once every fixer and verifier is judged.

## wording

the agent's choice in its plan, approved by the orchestrator acting as arbiter at plan approval

## Links

### In

- decided-by [review-fix-leaf](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-leaf.md)
- decided-by [review-fix-change-set](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-change-set.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
