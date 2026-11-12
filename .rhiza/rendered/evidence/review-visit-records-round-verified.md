# review-visit-records-round-verified

evidence · `01M76HKG4QXWAXZNSSW9QZ5WXM` · actor [claude-code-441c66c4](../actor/claude-code-441c66c4.md)

## result

PASS

## statement

Verified at 89df994ce (attempt 2, fixing attempt 1's findings f1/f2). go test -trimpath ./internal/api/... -run TestReviewVisitNeedsRound: ok (round-recorded and session done refuse naming rhiza run findings with no round; pass after an empty round; an earlier visit's round does not count). go test -trimpath ./cmd/rhiza -run TestScripts/plan-review-visit-round (GOTMPDIR outside the repo): PASS; it drives the plan workflow with fakes: a reviewer with no round has done refused (exit 3, naming rhiza run findings), the run fails at review and never reaches the critic; reviewers that record a round (a finding on visit 1, an empty round on visit 2 whose earlier round did not count) say done exit 0 and pass to the critic. internal/run and internal/workflow tests pass; describe-golden passes. Diff read: plan.toml and plan-milestone.toml review until = rhiza run round-recorded; isPlanReviewer excludes check/role/leaf sessions, so review.toml checks steps are unaffected. Note: TestScripts/plan-workflow and plan-milestone-backlog fail at HEAD but pass at this leaf's commit 71596b673; they break with sibling leaf review-check-refuses-stale-round (6718a349d, critic now goes to the new fresh step, expected 'then check'), not with this leaf.

## Links

### Out

- verifies [review-visit-records-round](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-visit-records-round.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
