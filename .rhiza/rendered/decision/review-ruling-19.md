# review-ruling-19

decision · `01M602SX94NGTNA2H4A577NEVH`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

The change-set fixer is journalled with its own review-fix and review-fix-judged entries, one per cycle, not as a leaf or check session. It gets the confirmed findings that name no leaf or a leaf outside the change set, and its commit carries a Review: trailer.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### In

- decided-by [review-fix-change-set](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-change-set.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
