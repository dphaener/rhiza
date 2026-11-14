# review-check-declared-checks-follows-main

decision · `01M7BZC6DPZ2BH0A3V8FE1M1EQ`

## reasoning

origin/main internal/api/run_review.go conductorReviewCheck (landed by critic-on-major, change set 01M79KGZ) already reads each declared check's latest round, fails on standing blockers and majors, and deliberately fails a declared check with no round so the conductor cannot skip a check and still pass its gate. That missing-round rule is landed, documented in the code and serves this change set's intent (every red hands back to conduct), so loosening it to 'passes in every other case' would weaken the gate and contradict main. The leaf's own core behaviour duplicates landed work, so the right settlement is to align the leaf with main and keep its build to the test and gaps main lacks (main has TestReviewCheckConductor but no TestReviewCheckDeclaredChecks), not to change the code or drop the leaf, which builtin-conductor-default depends on.

## statement

review-check-declared-checks restates main's landed conductorReviewCheck rather than re-specifying it: its statement and acceptance must match main, so review-check on a run past a step declaring checks fails (exit 1) when any declared check recorded no round (naming the missing checks) as well as when any latest round has a blocker or major the critic has not downgraded or dropped (naming each); it passes only when every declared check has a round and none has a standing blocker or major. The leaf's build is limited to what main lacks (the acceptance test covering minor-only pass, a named confirmed major failing, a missing round failing, and a downgrade or drop passing, plus any wiring the conductor default needs); it changes no landed behaviour.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f2](../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
