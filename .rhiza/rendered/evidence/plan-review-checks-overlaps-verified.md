# plan-review-checks-overlaps-verified

evidence · `01M6ADPZ7X9S4TQJRMDN20M1RQ` · actor [claude-code-8f5567a5](../actor/claude-code-8f5567a5.md)

## result

PASS

## statement

At 3a7069c4 the diff adds to internal/guide/text/plan-review.md a paragraph having the reviewer run `rhiza overlaps --change <cs> --json` for each plan change set and record each contradicting overlap as a `contradiction` finding naming the other change set; TestPlanReviewText in internal/guide/plan_test.go now asserts the phrases "rhiza overlaps --change <cs> --json" and "naming the other change set" beside "contradiction". go test ./internal/guide -count=1 -v: TestPlanReviewText, TestCaps and TestCapCatchesOversize PASS (plan-review 566 tokens, cap 800). make fmt-check vet test pass for every package. `rhiza overlaps --change` is accepted (overlaps.txtar exercises it). make self-check fails only on unrendered run actor records and index.md, which come from run entries outside this commit.

## Links

### Out

- verifies [plan-review-checks-overlaps](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-guide/plan-review-checks-overlaps.md)

## History

- `01M5SWTMR4G5EFRVRQHSEKVA15` Genesis reviews each plan against the model and the other runs in flight before work starts, as an agent check on the plan rather than the code.
