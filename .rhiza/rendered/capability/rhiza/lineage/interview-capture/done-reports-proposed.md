# done-reports-proposed

capability · `01M65071DAX3X3K3AEBNHAV49N` · parent [interview-capture](../interview-capture.md)

## acceptance

The done testscript's --json output lists both answer handles under decisions, and describe-golden is regenerated for the new response field

## statement

rhiza session done --json lists the handles of the Decisions it proposed under decisions, so the planner and the run see what was captured

## Links

### Out

- decided-by [done-reports-decisions](../../../../decision/done-reports-decisions.md)

### In

- verifies [done-reports-proposed-verified](../../../../evidence/done-reports-proposed-verified.md)
- verifies [done-reports-proposed-qa-1](../../../../evidence/done-reports-proposed-qa-1.md)
- verifies [pr-101-plan-captures-on-done](../../../../evidence/pr-101-plan-captures-on-done.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
