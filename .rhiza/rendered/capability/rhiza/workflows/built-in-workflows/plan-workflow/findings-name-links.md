# findings-name-links

capability · `01M76F8S09V3FZJC4T6RJHT5PY` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A test in internal/api (run_review_test.go) and the describe-golden testscript if a type changes: rhiza run findings with a finding whose node field is a depends-on link's ID between two change sets, and one naming a model link's ID, records the round; run status --json shows each with the link's type and ends; a verdict confirms one; a round-3 review-check adds a question naming the link with no blocks link. An ID that is neither a node nor a link still exits 2 with nothing written. Run: go test -trimpath ./internal/api -run <the new tests>, and go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden when regenerated.

## statement

rhiza run findings accepts a link ID (a model link, or a link between change sets such as depends-on) as a finding's target, journals the finding naming the link with its type and ends, and rhiza run status, rhiza run verdict and rhiza run review-check (whose round-3 question for a link finding names the link and links blocks to no node) handle it like a node finding.

## Links

### Out

- decided-by [link-findings-allowed](../../../../../decision/link-findings-allowed.md)
- derived-from [plan-review-round-unrecorded-plan](../../../../../source/plan-review-round-unrecorded-plan.md)

### In

- verifies [findings-name-links-verified](../../../../../evidence/findings-name-links-verified.md)
- verifies [findings-name-links-qa-1](../../../../../evidence/findings-name-links-qa-1.md)
- verifies [pr-452-fix-a-plan-review-round-whose](../../../../../evidence/pr-452-fix-a-plan-review-round-whose.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
