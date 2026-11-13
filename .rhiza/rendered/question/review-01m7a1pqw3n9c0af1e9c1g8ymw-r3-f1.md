# review-01m7a1pqw3n9c0af1e9c1g8ymw-r3-f1

question · `01M7A3K79BFD60Y3C3PDJJC0CS`

## statement

Plan review round 3 left finding f1 (oversized-leaf on daemon-plans-uncovered-proposals): The leaf carries two concerns: the shared covered predicate uncoveredProposals (new internal/api/proposals.go and proposals_test.go over journals of workflow plan, plan-milestone and plan-review, gated, failed, abandoned, live, pending, cross-change-set plan-milestone gating), which inbox-lists-uncovered-proposals also consumes, and the daemon's start wiring in daemonPoll1 (internal/api/daemon.go, 2224 lines, plus daemon_fold.go and daemon_projects.go: pending-start tracking, slot order after resumes and before ship picks, pause, low disk, load, memory and usage holds). Its acceptance is a new daemon testscript with about eleven scenarios plus the predicate's unit tests across six files, too much for one fresh session. The predicate is its own leaf that both the daemon leaf and the inbox leaf need. The critic confirmed it: Verified against the leaf 01M7A26C4: its statement bundles two concerns, the shared covered predicate uncoveredProposals (new internal/api/proposals.go over plan, plan-milestone and plan-review journals: live, pending, gated across change sets, failed, abandoned), which the inbox leaf inbox-lists-uncovered-proposals also consumes, and the daemon's start wiring in daemonPoll1 (internal/api/daemon.go, 2224 lines), where the new starts must slot between the resumes and daemonQueuedStarts/daemonPlanRequests (daemon.go:942-1168) with pending tracking (daemonHeld) and every hold. The acceptance is a new daemon testscript of about eleven scenarios plus the predicate's unit tests, touching six files. That is more than one concern and too much for one fresh session, and the predicate is a natural seam that both the daemon leaf and the inbox leaf need. Major oversized-leaf stands. How should it be settled?

## Links

### Out

- blocks [daemon-plans-uncovered-proposals](../capability/daemon-plans-uncovered-proposals.md)

### In

- answers [split-uncovered-proposals-predicate](../decision/split-uncovered-proposals-predicate.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
