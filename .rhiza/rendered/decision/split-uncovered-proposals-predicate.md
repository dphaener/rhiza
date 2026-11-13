# split-uncovered-proposals-predicate

decision · `01M7A3MPY2R9N56GAVM4HBFZB1`

## reasoning

The finding was confirmed by the critic: the leaf bundles a pure predicate over run journals with daemon start wiring in a 2224-line daemon.go, across six files and about eleven testscript scenarios plus unit tests, more than one fresh builder session. The predicate is a natural seam: it is pure, unit-testable on its own, and both the daemon leaf and the inbox leaf need it, so building it once as its own leaf removes the hidden ordering between those two leaves. Splitting along that seam keeps uncovered-proposals-get-plan-runs, failed-plan-run-covers and plan-review-skips-interview intact and changes no behaviour, only how the work is cut. Keeping the leaf whole or dropping scenarios would leave a major finding standing or weaken the acceptance the earlier Decisions require.

## statement

Split daemon-plans-uncovered-proposals in two. A new leaf owns the shared covered predicate uncoveredProposals (new internal/api/proposals.go and proposals_test.go): which proposed change sets with leaves no run covers, by run change set, workflow and journal (live, pending start, finished plan, plan-milestone or plan-review gated runs including cross-change-set plan-milestone gating, failed and not abandoned covers, abandoned covers nothing), accepted by its unit tests. daemon-plans-uncovered-proposals keeps only the daemon's start wiring in daemonPoll1 (slot order after resumes and before queued starts and plan requests, pending-start tracking, pause, low disk, load, memory and usage holds) with its daemon testscript, and inbox-lists-uncovered-proposals consumes the predicate; both depend on the predicate leaf. The covered rules and the scenarios themselves do not change.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7a1pqw3n9c0af1e9c1g8ymw-r3-f1](../question/review-01m7a1pqw3n9c0af1e9c1g8ymw-r3-f1.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
