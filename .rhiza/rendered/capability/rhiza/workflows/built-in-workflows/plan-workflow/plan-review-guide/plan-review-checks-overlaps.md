# plan-review-checks-overlaps

capability · `01M69DHPCMJZPPZF9076CWYS04` · parent [plan-review-guide](../plan-review-guide.md)

## acceptance

A case in internal/guide/plan_test.go asserts the plan review guide names rhiza overlaps and the contradiction finding, and the guide size cap test still passes

## statement

The plan review guide has the reviewer run rhiza overlaps --change <cs> --json for each plan change set and record each overlap that contradicts the plan as a contradiction finding naming the other change set

## Links

### Out

- decided-by [plan-review-overlaps-as-contradiction](../../../../../../decision/plan-review-overlaps-as-contradiction.md)

### In

- verifies [plan-review-checks-overlaps-verified](../../../../../../evidence/plan-review-checks-overlaps-verified.md)
- verifies [plan-review-checks-overlaps-qa-1](../../../../../../evidence/plan-review-checks-overlaps-qa-1.md)
- verifies [pr-151-genesis-reviews-each-plan](../../../../../../evidence/pr-151-genesis-reviews-each-plan.md)

## History

- `01M5SWTMR4G5EFRVRQHSEKVA15` Genesis reviews each plan against the model and the other runs in flight before work starts, as an agent check on the plan rather than the code.
