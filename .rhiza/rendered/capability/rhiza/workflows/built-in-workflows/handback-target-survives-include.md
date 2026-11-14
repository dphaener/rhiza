# handback-target-survives-include

capability · `01M7BYHZXGY011H2GCS1CZXX49` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/workflow -run TestIncludeHandbackTarget passes: a workflow including one whose step has on.fail = "handback:conduct" parses with that included step's on.fail = "handback:<workflow>/conduct", and the included conduct step present as <workflow>/conduct; an included handback:nosuch is refused with an error naming nosuch and the included workflow; done and failed targets still map as today. go test -trimpath ./internal/run -run TestHandbackIncluded passes: a run on such an including workflow whose included step fails hands back to <workflow>/conduct and continues from it. Touches internal/workflow/include.go, internal/workflow/include_test.go and internal/run/include_test.go.

## statement

A handback:<step> route inside an included workflow keeps pointing at the included step: when a workflow step includes another, include.go's inlined.outer maps handback:<step> to handback:<workflow>/<step> (as it prefixes a plain step target), instead of prefixing the whole target into <workflow>/handback:<step>; a handback naming a step the included workflow does not have is refused at parse, naming the step and the workflow.

## Links

### Out

- depends-on [run-handback](../workflow-runs/run-handback.md)
- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)

### In

- depends-on [rhiza-repo-conductor-default](../workflow-overrides/rhiza-repo-conductor-default.md)
- blocks [review-01m7ah6rayrzax6z6r63q0m2qk-r6-f3](../../../../question/review-01m7ah6rayrzax6z6r63q0m2qk-r6-f3.md)
- verifies [handback-target-survives-include-verified](../../../../evidence/handback-target-survives-include-verified.md)
- verifies [handback-target-survives-include-qa-2](../../../../evidence/handback-target-survives-include-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
