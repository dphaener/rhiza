# plan-review-workflow-verified

evidence · `01M7AB9JJSRETXTBJGW0B90FD1` · actor [claude-code-d9e542ac](../actor/claude-code-d9e542ac.md)

## result

PASS

## statement

At 208a5fb re-ran the Checks block: go test -trimpath ./internal/workflow/... ok (TestBuiltinPlanReview: start plan/milestone-check, no interview or planner step, gate present); go test -trimpath ./cmd/rhiza -run TestScripts/plan-review-workflow ok (workflow describe lists no interview; run start --workflow plan-review --autonomous reaches plan/gate via milestone-check and review, only the reviewer prompt opens, change set accepted under rule autonomous); go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-ls|plan-workflow|health|plan-milestone)' ok; rhiza ship generate exit 0; rhiza render exit 0. Diff touches only plan-review.toml (include plan, skip interview), workflow_test.go and the new txtar plus the actor record; no stubs or TODOs; fakes only for rex, claude and gh.

## Links

### Out

- verifies [plan-review-workflow](../capability/plan-review-workflow.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
