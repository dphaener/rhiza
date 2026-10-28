# review-ruling-17

decision · `01M602SWEGN0888T75ZZ6XD4P7`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

Leaf fixers and the verifiers that re-run after them are journalled as leaf-session and leaf-attempt entries on the checks step, with Attempt equal to the cycle, so the leaves step's start, poll, judge and run findings --leaf paths are reused. No leaf outcome entry is written on a checks step.

## wording

the agent's choice in its plan, approved by the orchestrator acting as arbiter at plan approval

## Links

### In

- decided-by [review-fix-leaf](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-leaf.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
