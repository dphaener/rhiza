# daemon-plans-uncovered-proposals-verified-2

evidence · `01M7ASNHSSN0A83FRM6T1JWGCM` · actor [claude-code-9400419f](../actor/claude-code-9400419f.md)

## result

PASS

## statement

Attempt 2 at 514c416f2: go vet ./internal/api/ exit 0; go test -trimpath ./internal/api/ -run 'TestUncoveredProposals|TestDaemonPlansProposals|TestDaemonHeld|Unjournaled|TestInbox' exit 0; go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-' exit 0 (daemon-plans-proposals.txtar included); go test -trimpath ./internal/api/... hit the default 10m package timeout twice under host load ~13 (running TestReadBudgetApply, untouched), and with -timeout 40m -v passed every test (ok, 638s; slowest TestReadBudget 141s and TestLinkEndpointKinds 129s, unrelated). The 'before the fix: exit 1' check records the pre-fix state and is not re-runnable at HEAD. Diff read: uncoveredProposals now covers only on live/held, a plan/plan-milestone/plan-review gated entry, or a failed non-abandoned run, matching the statement; daemonHeld keeps unjournaled starts of proposed as well as accepted change sets in daemonPoll1 and daemonDemand; TestUncoveredProposals updated (succeeded plan-review without gate covers nothing) and TestDaemonPlansProposalsPending polls past grace and hold. Only internal/api/daemon.go, daemon_projects.go, proposals.go, proposals_test.go plus .rhiza records changed; no TODOs or stubs; fakes only.

## Links

### Out

- verifies [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
