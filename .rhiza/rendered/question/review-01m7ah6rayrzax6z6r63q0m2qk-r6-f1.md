# review-01m7ah6rayrzax6z6r63q0m2qk-r6-f1

question · `01M7BZ7E6SHBJCMR8SH6C2MWHX`

## statement

Plan review round 6 left finding f1 (contradiction on agent-step-checks): Already built on main by change set 01M79KGZ (commit c09395b12, #575, leaf review-three-checks): internal/api/run_review.go reviewStep already treats an agent step with Checks as a review step, reviewCheck refuses an undeclared check naming the step's checks, and internal/workflow/workflow_test.go has TestAgentStepChecks. Main also declares them as [[step.conduct.check]] tables with a name key, not the checks = ["plan", "code", "qa"] key this leaf specifies. The plan should drop the leaf or restate it as the gap that is left, in main's syntax. The critic confirmed it: Verified on origin/main: c09395b12 (#575) is on main, internal/workflow/workflow_test.go:2808 has TestAgentStepChecks, and main declares checks as [[step.conduct.check]] tables with name = "..." (run_review_test.go:302, daemon_picks_test.go:231), not checks = ["plan","code","qa"]. The leaf re-specifies landed behaviour in a syntax main's parser does not use, so a builder would contradict the code. Major stands. How should it be settled?

## Links

### Out

- blocks [agent-step-checks](../capability/rhiza/workflows/built-in-workflows/built-in-review/agent-step-checks.md)

### In

- answers [agent-step-checks-already-built](../decision/agent-step-checks-already-built.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
