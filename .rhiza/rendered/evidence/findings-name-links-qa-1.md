# findings-name-links-qa-1

evidence · `01M76JXRVF4JESXDGEY9TGZFP9` · actor [claude-code-3f408ce1](../actor/claude-code-3f408ce1.md)

## result

PASS

## statement

QA at c20429608, built binary (go build -trimpath) in a scratch repo with a hand-written plan run journal (builtin plan.toml definition), RHIZA_RUN set: rhiza run findings with a change-set depends-on link ID (also by 12-char prefix), a model depends-on link ID and a node recorded the round (exit 0); text and run status --json showed link, link_type depends-on, link_from and link_to; run verdict confirmed the link finding; round-3 run review-check added question review-…-r3-f1 naming 'link depends-on <id> from <a> to <b>' with links_in and links_out empty; an ID that is neither node nor link exited 2 with the journal unchanged.

## Links

### Out

- verifies [findings-name-links](../capability/rhiza/workflows/built-in-workflows/plan-workflow/findings-name-links.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
