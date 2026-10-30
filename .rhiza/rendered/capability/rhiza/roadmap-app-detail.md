# roadmap-app-detail

capability · `01M6734AZGT872HRJ9ZRBWAABE` · parent [rhiza](../rhiza.md)

## acceptance

A test builds the app over fake change.show, ready and run results, presses enter on a change set row and asserts the detail lists its intent, each leaf's statement and acceptance, each blocking question's statement and the latest run's state and step; esc returns to the list with the same row selected

## statement

Enter on a change set in the roadmap app opens its detail: intent, kind and state, its leaves each with statement and acceptance, the open questions that block it, and its latest run with the run's state and current step; esc goes back to the list with the cursor where it was

## Links

### Out

- serves [m5-roadmap-actionable](../../intent/milestone-5/m5-roadmap-actionable.md)
- decided-by [m5-roadmap-inbox-first](../../decision/m5-roadmap-inbox-first.md)

### In

- verifies [roadmap-app-detail-verified](../../evidence/roadmap-app-detail-verified.md)
- verifies [roadmap-app-detail-qa-1](../../evidence/roadmap-app-detail-qa-1.md)
- verifies [pr-120-actionable-roadmap-app](../../evidence/pr-120-actionable-roadmap-app.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
