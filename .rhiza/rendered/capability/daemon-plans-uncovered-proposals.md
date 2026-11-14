# daemon-plans-uncovered-proposals

capability · `01M7A26C4ZRFSRVS8ND06BE7NJ`

## acceptance

A new cmd/rhiza/testdata/script/daemon-plans-proposals.txtar (Rex faked, RHIZA_TEST_DETACH_LOG): a change set proposed with a leaf and no run makes rhiza daemon --once start one plan-review run naming it with --autonomous and the daemon's permission mode; a second poll while that start is pending starts nothing more; a proposed change set with no leaves, one a finished plan run's gated entry names, one a finished plan-milestone run on another change set (the record change set) left proposed and whose gated entry names it, one a finished plan-review run's gated entry names (started --workflow plan-review without --new, gate left it proposed under rhiza.autonomy.accept owner), one whose plan-review run ended failed and one whose plan run ended failed start nothing on this or a following poll; one whose only run was abandoned gets a start; with --max-runs 1 and a ready accepted change set the proposal's start and the ship pick share the slot without exceeding it; under rhiza daemon pause nothing starts. internal/api/proposals_test.go covers uncoveredProposals over journals of workflow plan, plan-milestone and plan-review, gated, failed, abandoned and live, including a plan-milestone run whose own change set differs from the gated change set it names. go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-' pass. Touches internal/api/daemon.go, internal/api/daemon_fold.go, internal/api/daemon_projects.go, internal/api/proposals.go, internal/api/proposals_test.go, cmd/rhiza/testdata/script/daemon-plans-proposals.txtar

## statement

Each daemon poll starts, as a detached rhiza run start --foreground <change> --workflow plan-review --autonomous --permission-mode <the daemon's mode>, a plan-review run for every proposed change set with leaves that no run covers. Covered is defined by a run's change set, workflow and journal, never by run.State.IsPlanRun or isPlanWorkflow (which know only plan and plan-milestone, and stay unchanged): a change set is covered while any run on it is live or its start is pending and not yet journaled, when any finished run whose workflow is plan, plan-milestone or plan-review (or that run.State.IsPlanRun reports), whatever its own change set (a plan-milestone run on the record change set gates every change set serving the intents it creates), has a latest gated entry naming it, or when any run on it, plan-review included, ended failed and was not abandoned. These starts count against max-runs after the resumes and before the ship picks, and stop under pause, low disk, the load and memory holds and a usage hold like any other start; one shared predicate (uncoveredProposals in internal/api) decides covered, for the daemon and the inbox

## Links

### Out

- depends-on [plan-review-workflow](plan-review-workflow.md)
- decided-by [uncovered-proposals-get-plan-runs](../decision/uncovered-proposals-get-plan-runs.md)
- decided-by [failed-plan-run-covers](../decision/failed-plan-run-covers.md)

### In

- depends-on [inbox-lists-uncovered-proposals](inbox-lists-uncovered-proposals.md)
- blocks [review-01m7a1pqw3n9c0af1e9c1g8ymw-r3-f1](../question/review-01m7a1pqw3n9c0af1e9c1g8ymw-r3-f1.md)
- verifies [daemon-plans-uncovered-proposals-verified](../evidence/daemon-plans-uncovered-proposals-verified.md)
- verifies [daemon-plans-uncovered-proposals-qa-1](../evidence/daemon-plans-uncovered-proposals-qa-1.md)
- verifies [daemon-plans-uncovered-proposals-verified-2](../evidence/daemon-plans-uncovered-proposals-verified-2.md)
- verifies [daemon-plans-uncovered-proposals-verified-5](../evidence/daemon-plans-uncovered-proposals-verified-5.md)
- verifies [pr-559-fix-a-change-set-proposed](../evidence/pr-559-fix-a-change-set-proposed.md)

## History

- `01M7A1PXPS84ZD222S9ADHF0SV` Fix A change set proposed outside a plan run is never accepted and nothing says so. Milestone 36's change sets 1-4 (01M79KGHYF, 01M79KGRYE, 01M79KGZ51, 01M79KH59D) were proposed by PR #528 (merged 05:57Z on 2026-11-13) and milestone 37's 2-6 (01M79NKERT, 01M79NM3RJ, 01M79NMNGD, 01M79NN43P, 01M79NNP2B) by PR #531 (06:00Z), each with leaves written with the owner. With rhiza.autonomy.accept = auto they sat proposed for about 3 hours while slots were free: acceptance only happens at a plan run's gate (runRunGate, internal/api/run_review.go), the daemon starts no plan run for a change set that is already proposed with leaves, rhiza ready lists them under Not ready as 'not accepted: proposed', and the inbox says nothing waits. Unstuck by the orchestrator at 08:55Z: rhiza run start <cs> --workflow plan --autonomous --permission-mode auto for each of the nine. Wanted: the daemon starts a plan (review) run for a proposed change set that no plan run covers, as it does for a draft with no leaves; or the inbox lists such proposals when accept is auto. Related, different: 01M6PQTF5E (owner-run plan gate leaves proposals).
