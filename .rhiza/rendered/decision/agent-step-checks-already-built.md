# agent-step-checks-already-built

decision · `01M7BZ9EVT7XKQWCRYE55NGE3H`

## reasoning

Verified on origin/main: c09395b12 (#575) is merged; internal/workflow/workflow_test.go:2808 TestAgentStepChecks parses an agent step's checks in file order and refuses a check with no name and a duplicate name; internal/api/run_review.go:205 accepts --check at an agent step with Checks and refuses otherwise; internal/api/run_review_test.go has a test recording findings and verdicts under --check plan at the conduct step and refusing --check scope with a message naming the declared checks, plus TestReviewCheckConductor. Every acceptance point of the leaf is met by landed code, and its checks = [...] syntax contradicts the parser main uses, so a builder would either redo landed work or break it. Restating it as a remaining gap would yield an empty leaf, so dropping it is the answer that best fits the code and the change set's intent (promote the conductor default, not re-add check declarations).

## statement

Drop the agent-step-checks leaf from this change set: its behaviour is already on main (change set 01M79KGZ, #575), declared as [[step.conduct.check]] tables with a name key, and no gap is left to build. Dependents (review-check-declared-checks, builtin-conductor-default) rely on main's existing support and main's [[step.<name>.check]] name = "..." syntax, not checks = [...].

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f1](../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f1.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
