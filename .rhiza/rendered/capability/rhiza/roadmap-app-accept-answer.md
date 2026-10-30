# roadmap-app-accept-answer

capability · `01M6734EA3V7228YH8DGYVYDQA` · parent [rhiza](../rhiza.md)

## acceptance

A test over a temporary repository with a proposed change set presses a, sees the what-if, presses y and asserts the change set is accepted; another with a change set blocked by a question opens its detail, picks the question, types an answer and enter, and asserts a Decision quoting the typed text answers the question, rhiza ready reports the change set ready, and rhiza check passes

## statement

On the selected change set, a accepts a proposed one after showing the accept what-if and a y/n; in the detail, choosing a blocking question and typing an answer records a Decision with the owner's typed words as its quote, linked answers to the question, in that change set, so the change set becomes ready; the list reloads after either

## Links

### Out

- serves [m5-roadmap-actionable](../../intent/milestone-5/m5-roadmap-actionable.md)
- decided-by [m5-roadmap-inbox-first](../../decision/m5-roadmap-inbox-first.md)

### In

- verifies [roadmap-app-accept-answer-verified](../../evidence/roadmap-app-accept-answer-verified.md)
- verifies [roadmap-app-accept-answer-qa-1](../../evidence/roadmap-app-accept-answer-qa-1.md)
- verifies [roadmap-app-accept-answer-verified-2](../../evidence/roadmap-app-accept-answer-verified-2.md)
- verifies [pr-120-actionable-roadmap-app](../../evidence/pr-120-actionable-roadmap-app.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
