# plan-review-checks-overlaps-qa-1

evidence · `01M6AEBTD7GNNZNK9EJF2KHVRW` · actor [claude-code-e44db580](../actor/claude-code-e44db580.md)

## result

PASS

## statement

At 04deb2ed the built binary's rhiza guide plan-review prints, at line 24, the paragraph telling the reviewer to run `rhiza overlaps --change <cs> --json` for each plan change set and record each contradicting overlap as a `contradiction` finding naming the other change set. I ran that command as the guide says against a scratch repo and it returned the overlap JSON the guide describes. go test ./internal/guide -count=1 passes (TestPlanReviewText and the cap tests).

## Links

### Out

- verifies [plan-review-checks-overlaps](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-guide/plan-review-checks-overlaps.md)

## History

- `01M5SWTMR4G5EFRVRQHSEKVA15` Genesis reviews each plan against the model and the other runs in flight before work starts, as an agent check on the plan rather than the code.
