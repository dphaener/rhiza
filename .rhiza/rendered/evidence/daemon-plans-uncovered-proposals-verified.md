# daemon-plans-uncovered-proposals-verified

evidence · `01M7AEENQJ2X7C1TK5R7Q01QMR` · actor [claude-code-364a199b](../actor/claude-code-364a199b.md)

## result

PASS

## statement

At c4033e93b re-ran the Checks block: go test -trimpath ./internal/api -run 'Daemon|Proposals|Usage|Dispatch|Inbox|Idle|Guide' ok (251s); go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-' ok (incl. new daemon-plans-proposals.txtar); rhiza ship generate exit 0, 0 files changed; rhiza render left the tree clean. TestUncoveredProposals (plan, plan-milestone on the record change set, plan-review gated, failed, abandoned, live, held) and TestDaemonPlansProposalsPending (pending start held across polls in one process, pause, max-runs shared with the ship pick) pass. The acceptance's bare go test -trimpath ./internal/api/... hit Go's default 10m timeout under host load 8 in the untouched TestLinkEndpointKinds (100s alone); the repo's make test uses -timeout 20m. gofmt and vet are clean, with no TODOs. The diff matches the statement: one shared uncoveredProposals predicate, starts after the resumes and before the ship picks, counted in extra and daemonDemand, skipped when stopped or paused. daemon_fold_test.go, outside the listed files, only adds failed plan-review runs to the fixture so existing tests see no plan-review starts.

## Links

### Out

- verifies [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
