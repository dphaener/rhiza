# inbox-lists-uncovered-proposals-verified

evidence · `01M7AN1720B0FY01YY0A9WJ9WT` · actor [claude-code-2addd897](../actor/claude-code-2addd897.md)

## result

PASS

## statement

At d4f3f9e1b (attempt 2): go test -trimpath ./internal/api/... -skip Budget -timeout 25m ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-|ready)' ok (includes inbox-uncovered-proposal.txtar and inbox-owner-run-accept.txtar); rhiza ship generate and rhiza render exit 0 with no diff. Diff read: inboxUncoveredProposals/uncoveredAcceptItems add one accept item per proposed change set with leaves that uncoveredProposals leaves uncovered, detail 'waiting for acceptance: no plan run covers it', next ownerAcceptNext, regardless of rhiza.autonomy.accept, skipping change sets already holding a gated accept item. The txtar asserts 3 items under accept auto and owner, 2 after a failed plan run, 1 after accept, 0 with a live plan run; inbox-owner-run-accept.txtar gives the gated change set a leaf and asserts exactly one accept item; inbox_test.go covers live/failed/abandoned runs, gated dedup and other kinds. Attempt 1 findings f1 and f2 are addressed. ready code is untouched. No stubs or TODOs.

## Links

### Out

- verifies [inbox-lists-uncovered-proposals](../capability/inbox-lists-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
