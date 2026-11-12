# review-check-refuses-stale-round-verified

evidence · `01M76HZ9Z9NQY3SEEK6QDBZ3XR` · actor [claude-code-5d4ce6fe](../actor/claude-code-5d4ce6fe.md)

## result

PASS

## statement

Verifier at 6718a349d: go test -trimpath ./internal/api ./internal/workflow/... -run 'TestReviewCheckRefusesStaleRound|TestPlanFreshStep|TestPlanBuiltin' passes; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes (golden has the fresh flag); full internal/workflow and internal/run packages pass; full internal/api passes at this commit from a scratch worktree outside ~/code with TestReadBudget skipped (a timing budget this change does not touch, 190s under host load); run inside ~/code, TestDaemonSkipsReadyWithUnresolvedWorkflow fails on the shared ~/code/.rhiza-root usage.live.json changing, which is environmental (it passes from the scratch worktree). Read the diff: State.RoundStale (rounds <= rounds at the latest revise entry) gates review-check before judging, with a Conflict naming round N as predating the latest revise and writing nothing; --fresh returns pass after that gate only. The test covers round 2 and round 3 stale for both forms with no asked entry, and a later fresh round passing both. plan.toml and plan-milestone.toml add step fresh (run = rhiza run review-check --fresh, on.pass = check, on.fail = review) and critic on.pass = fresh, checked by TestPlanFreshStep. An extra overlay test (no tree edit) confirmed: no revise -> --fresh passes and check judges round 1 as before; round 2 then revise -> stale refusal naming round 2; round 3 after revise -> --fresh exits 0 with no asked entry and review-check passes the stale gate into the unchanged judging path. Note: the leaf's own test makes its fresh round round 4 with no findings rather than a confirmed round 3.

## Links

### Out

- verifies [review-check-refuses-stale-round](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-check-refuses-stale-round.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
