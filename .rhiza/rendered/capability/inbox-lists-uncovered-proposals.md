# inbox-lists-uncovered-proposals

capability · `01M7A26CAT1FTRENTDNPR4XD18`

## acceptance

A new cmd/rhiza/testdata/script/inbox-uncovered-proposal.txtar: a change set proposed by hand with a leaf and no run shows in rhiza inbox --json as one accept item with that detail and next, under accept auto and under accept owner; once a plan run on it is live (or ended failed) the accept item is gone; once accepted it is gone; inbox-owner-run-accept.txtar still lists a gated proposal once. go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-|ready)' pass. Touches internal/api/inbox.go, internal/api/inbox_test.go, cmd/rhiza/testdata/script/inbox-uncovered-proposal.txtar

## statement

rhiza inbox lists every proposed change set that no plan run covers (uncoveredProposals) as an accept item, detail 'waiting for acceptance: no plan run covers it', next rhiza accept --rule owner --change <id>, whatever rhiza.autonomy.accept says; a proposal a gate left for the owner keeps its existing gated accept item, and no change set is listed twice; rhiza ready keeps its 'not accepted: proposed' reason

## Links

### Out

- depends-on [daemon-plans-uncovered-proposals](daemon-plans-uncovered-proposals.md)
- decided-by [uncovered-proposals-get-plan-runs](../decision/uncovered-proposals-get-plan-runs.md)
- decided-by [failed-plan-run-covers](../decision/failed-plan-run-covers.md)

### In

- verifies [inbox-lists-uncovered-proposals-verified](../evidence/inbox-lists-uncovered-proposals-verified.md)
- verifies [inbox-lists-uncovered-proposals-qa-1](../evidence/inbox-lists-uncovered-proposals-qa-1.md)
- verifies [pr-559-fix-a-change-set-proposed](../evidence/pr-559-fix-a-change-set-proposed.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
