# plan-review-workflow

capability · `01M7A26C10VA8H52GYENGK7NC4`

## acceptance

go test -trimpath ./internal/workflow/... covers parsing plan-review (its first step after expansion is the included milestone-check, no interview step); a new cmd/rhiza/testdata/script/plan-review-workflow.txtar shows rhiza workflow describe plan-review listing no interview step and a rhiza run start --workflow plan-review --autonomous --no-wait on a proposed change set with a leaf, with the review and critic faked as the plan-workflow scripts fake them, reaching the gate, which accepts it under rule autonomous with rhiza.autonomy.accept auto; testscripts that list the builtin workflow names (go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-ls|plan-workflow)') pass, updated where they print the list. Touches internal/workflow/builtin/plan-review.toml, internal/workflow/workflow_test.go, cmd/rhiza/testdata/script/plan-review-workflow.txtar

## statement

A builtin workflow plan-review (internal/workflow/builtin/plan-review.toml) is one workflow step that includes plan with skip = ["interview"], so rhiza run start <change> --workflow plan-review on a proposed change set with leaves enters at milestone-check, then the plan review, then the gate, and never opens a planner session

## Links

### Out

- decided-by [plan-review-skips-interview](../decision/plan-review-skips-interview.md)
- decided-by [uncovered-proposals-get-plan-runs](../decision/uncovered-proposals-get-plan-runs.md)

### In

- depends-on [daemon-plans-uncovered-proposals](daemon-plans-uncovered-proposals.md)
- verifies [plan-review-workflow-verified](../evidence/plan-review-workflow-verified.md)
- verifies [plan-review-workflow-qa-1](../evidence/plan-review-workflow-qa-1.md)
- verifies [pr-559-fix-a-change-set-proposed](../evidence/pr-559-fix-a-change-set-proposed.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
