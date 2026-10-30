# roadmap-app-list

capability · `01M67349MAEY43NAK3P9PG0H74` · parent [rhiza](../rhiza.md)

## acceptance

A test in internal/tui builds the roadmap app over fake roadmap and ready results and asserts the view has the group headings, a row per change set with state, title and not-ready reasons, and no JSON braces; that j moves the cursor to the next change set row, skipping headings; and the --once test in cmd/rhiza/testdata/script still compares equal to rhiza roadmap --json

## statement

The roadmap app shows the roadmap as a cursor list, not JSON: a heading per milestone group and the backlog, one row per change set with its short ID, state, title, and the reasons it is not ready; j/k or the arrows move the cursor, r reloads, q quits, and rhiza tui app roadmap --once still prints what rhiza roadmap --json prints

## Links

### Out

- serves [m5-roadmap-actionable](../../intent/milestone-5/m5-roadmap-actionable.md)
- decided-by [m5-roadmap-inbox-first](../../decision/m5-roadmap-inbox-first.md)

### In

- verifies [roadmap-app-list-verified](../../evidence/roadmap-app-list-verified.md)
- verifies [roadmap-app-list-qa-1](../../evidence/roadmap-app-list-qa-1.md)
- verifies [pr-120-actionable-roadmap-app](../../evidence/pr-120-actionable-roadmap-app.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
