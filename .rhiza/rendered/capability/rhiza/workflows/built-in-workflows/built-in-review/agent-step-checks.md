# agent-step-checks

capability · `01M7AMS2SDWJ2NYPKKVW4HXB7G` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/workflow -run TestAgentStepChecks shows an agent step's checks parsed in file order and rhiza workflow check refusing an empty or repeated check name; go test -trimpath ./internal/api -run TestRunFindingsAgentStepChecks shows rhiza run findings --check and rhiza run verdict --check recording a round for a declared check at an agent step, refusing an undeclared check with a message naming the declared ones, and refusing --check at an agent step with no checks. Touches internal/workflow/workflow.go, internal/workflow/workflow_test.go, internal/api/run_review.go and internal/api/run_review_test.go.

## statement

An agent step may declare review checks by name (checks = ["plan", "code", "qa"] in its table), with no critic or fixer roles of its own. While a run is at that step, rhiza run findings --check and rhiza run verdict --check accept a check it declares and record the round as they do in a checks step, and refuse a check it does not declare with a message naming the declared ones; at an agent step that declares no checks, --check is refused as today. Checks steps work as they do now.

## Links

### Out

- decided-by [conductor-promote-builtin-now](../../../../../decision/conductor-promote-builtin-now.md)

### In

- depends-on [review-check-declared-checks](review-check-declared-checks.md)
- depends-on [builtin-conductor-default](../builtin-conductor-default.md)
- blocks [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f1](../../../../../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f1.md)
- verifies [agent-step-checks-verified](../../../../../evidence/agent-step-checks-verified.md)
- verifies [agent-step-checks-qa-2](../../../../../evidence/agent-step-checks-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
