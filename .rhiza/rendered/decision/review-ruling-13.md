# review-ruling-13

decision · `01M602STSDCAD6066F08ECKPPD`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

Check and critic sessions run in parallel, each in a throwaway detached checkout at the run worktree's HEAD, with RHIZA_ROOT pointing their rhiza writes at the run worktree. The checkout comes from its own OpenCheckCheckout path rather than the leaves' OpenCheckout, so a check's checkout can never collide with a leaf verifier's.

## wording

the agent's choice in its plan, approved by the orchestrator acting as arbiter at plan approval

## Links

### In

- decided-by [review-check-worktree](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-check-worktree.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
