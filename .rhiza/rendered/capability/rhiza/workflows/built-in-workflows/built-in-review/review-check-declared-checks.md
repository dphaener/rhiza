# review-check-declared-checks

capability · `01M7AMS4TR5TGACVFJEDFEFZYE` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/api -run TestReviewCheckDeclaredChecks shows, for a run at a later step than an agent step declaring plan, code and qa, review-check passing with only minor findings, failing with exit 1 and naming a confirmed major of the code check, and passing once the critic downgraded that major to minor or dropped it; go test -trimpath ./internal/api -run 'TestReviewCheckRefusesStaleRound|TestReviewCheckFindingsFlag' still passes, both driving review-check on a plan-review run (workflow plan). Touches internal/api/run_review.go and internal/api/run_review_test.go.

## statement

rhiza run review-check, on a run whose workflow has an agent step declaring checks and which is past that step (at the gate), reads the latest round of each declared check and fails (exit 1) when any has a blocker or major finding the critic has not downgraded or dropped, naming each such finding; otherwise it passes. It leaves no question and has no round-3 pass. On a plan-review run, review-check is unchanged.

## Links

### Out

- depends-on [agent-step-checks](agent-step-checks.md)
- depends-on [critic-on-major](../critic-on-major.md)
- decided-by [conductor-promote-builtin-now](../../../../../decision/conductor-promote-builtin-now.md)

### In

- depends-on [builtin-conductor-default](../builtin-conductor-default.md)
- blocks [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f2](../../../../../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f2.md)
- verifies [review-check-declared-checks-verified](../../../../../evidence/review-check-declared-checks-verified.md)
- verifies [review-check-declared-checks-qa-2](../../../../../evidence/review-check-declared-checks-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
