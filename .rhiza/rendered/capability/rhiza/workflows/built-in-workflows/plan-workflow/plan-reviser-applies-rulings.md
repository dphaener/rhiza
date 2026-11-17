# plan-reviser-applies-rulings

capability · `01M7MPK73E4BCQACH2KH8XSYH2` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestPlanReviserText' passes, asserting the rulings section names rhiza run status --json, the decided-by link from a live node to the Decision with its three cases (the finding's node, the from node of a link finding, each replacing node when the ruling retires the finding's node), --record for a landed record, and the rule never to edit the Decision, with the existing headings still in order.

## failure_waiver

Guide text only: the test asserts the text is present; there is no runtime failure path in this leaf (the run failing on an unapplied ruling is plan-workflow-applies-rulings' failure case).

## files

internal/guide/text/plan-reviser.md, internal/guide/plan_reviser_test.go

## must_not

Change the confirmed-findings instructions or other guide topics; tell the reviser to write a quote or a new Decision.

## statement

rhiza guide plan-reviser has a section for a visit told to apply the arbiter's rulings: read the run's questions with a decision in rhiza run status --json, and for each change the node of the finding it answers as the Decision says (with rhiza apply in the change set holding the node, --record for a landed record) and link a live node decided-by the Decision: the finding's node when it stays live; the link's from node when the finding targets a link; and each node that replaces it when the ruling retires the finding's node (an oversized-leaf split links every new leaf). Never edit, retire or re-link the Decision, and change nothing a ruling does not name.

## Links

### Out

- decided-by [plan-rulings-applied-before-gate](../../../../../decision/plan-rulings-applied-before-gate.md)
- decided-by [ruling-applied-means-decided-by](../../../../../decision/ruling-applied-means-decided-by.md)

### In

- verifies [plan-reviser-applies-rulings-verified](../../../../../evidence/plan-reviser-applies-rulings-verified.md)
- verifies [plan-reviser-applies-rulings-qa-1](../../../../../evidence/plan-reviser-applies-rulings-qa-1.md)
- verifies [plan-reviser-applies-rulings-qa-2](../../../../../evidence/plan-reviser-applies-rulings-qa-2.md)
- verifies [plan-reviser-applies-rulings-qa-3](../../../../../evidence/plan-reviser-applies-rulings-qa-3.md)
- verifies [plan-reviser-applies-rulings-qa-4](../../../../../evidence/plan-reviser-applies-rulings-qa-4.md)
- verifies [plan-reviser-applies-rulings-qa-5](../../../../../evidence/plan-reviser-applies-rulings-qa-5.md)
- verifies [plan-reviser-applies-rulings-qa-6](../../../../../evidence/plan-reviser-applies-rulings-qa-6.md)
- verifies [plan-reviser-applies-rulings-qa-7](../../../../../evidence/plan-reviser-applies-rulings-qa-7.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
