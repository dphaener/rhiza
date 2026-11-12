# findings-name-links-verified

evidence · `01M76J1NFTW2MYRMEZZGPSR37Q` · actor [claude-code-e13e6def](../actor/claude-code-e13e6def.md)

## result

PASS

## statement

At 897760da8, ran go test -trimpath -count=1 ./internal/api -run TestReviewFindingsNameLinks (PASS), ./internal/run/... (ok) and ./cmd/rhiza -run TestScripts/describe-golden plus the review/findings scripts (ok). I read the diff from 89df994ce. run findings resolves a link ID (in the run's change set or a plan change set's view) when node resolution fails, and journals link, link_type, link_from and link_to (run.Finding, api.RunFinding). The test records a depends-on link between two change sets and a model depends-on link. run.status returns each with type and ends, and its markdown names 'link depends-on <id> from <a> to <b>'. A verdict confirms each one. A round-3 review-check adds a question with no blocks link. An unknown ID exits 2 and journals no round. Attempt 1's two findings (no ends recorded, no run.status check) are fixed.

## Links

### Out

- verifies [findings-name-links](../capability/rhiza/workflows/built-in-workflows/plan-workflow/findings-name-links.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
