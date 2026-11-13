# daemon-plans-uncovered-proposals-verified-5

evidence · `01M7AYBDZJ34M61Q5WS56W8XQH` · actor [claude-code-25d088c3](../actor/claude-code-25d088c3.md)

## result

PASS

## statement

Attempt 5 at 309e097b3 (fixer commit: .rhiza log entries, rendered actor pages and index.md only; no Go or test code changed), in a clean detached worktree: Checks block re-run: make self-check (env -i) exit 0; rhiza ship generate exit 0, 0 files staged; rhiza render exit 0 (wrote only the verifier's own new actor page and index.md). Acceptance: go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-' ok (38.8s, daemon-plans-proposals.txtar included); go test -trimpath ./internal/api/... as written hit the default 10m package timeout twice on this 4-core host under load 6-9 (still running TestLinkEndpointKinds, which is untouched and progressing, not hung), the same as at attempt 2. With -timeout 40m every test passed except TestReadBudget, a CPU timing budget: driver poll idle median 6.249ms CPU over its 5ms budget, at 9% CPU pressure on the shared host. This leaf does not touch driver poll, and TestReadBudget passed at attempt 2 with the same leaf code. Both files named in the acceptance exist. Diff read: nothing outside .rhiza changed, no stubs or TODOs, fakes only.

## Links

### Out

- verifies [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
