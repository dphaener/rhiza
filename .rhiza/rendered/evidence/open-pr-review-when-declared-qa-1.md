# open-pr-review-when-declared-qa-1

evidence · `01M7REYTTJ2QWQT5G5R2RDMX75` · actor [claude-code-8a84713f](../actor/claude-code-8a84713f.md)

## result

PASS

## statement

QA at e288b6b36 with a scratch build: myteam-flow's push (ship open-pr) passed with no review journaled; the same workflow with a checks step added refused push with 'has no review journaled: its workflow declares a review'; ship open-pr outside a run refused. TestShipOpenPR|TestReviewJournaled|TestReviewDeclared, internal/workflow, run-ship-change and ship-open-pr scripts pass.

## Links

### Out

- verifies [open-pr-review-when-declared](../capability/ship-through-rhiza/open-pr-review-when-declared.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
