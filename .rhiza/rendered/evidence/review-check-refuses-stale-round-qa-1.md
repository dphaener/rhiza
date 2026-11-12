# review-check-refuses-stale-round-qa-1

evidence · `01M76JXS04BPC7JHBMBQ55GZKN` · actor [claude-code-3f408ce1](../actor/claude-code-3f408ce1.md)

## result

PASS

## statement

QA at c20429608, built binary (go build -trimpath) in a scratch repo with a hand-written plan run journal (builtin plan.toml definition), RHIZA_RUN set: after a revise visit with no later round, run review-check and run review-check --fresh both exited 3 naming round 1 as predating the latest revise and wrote nothing; with rounds journaled after each revise, --fresh exited 0 and review-check judged as before (round 2 fail exit 1, round 3 pass adding a question). workflow describe plan shows critic on.pass -> fresh, fresh on.pass -> check, on.fail -> review; plan-milestone.toml has the same.

## Links

### Out

- verifies [review-check-refuses-stale-round](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-check-refuses-stale-round.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
