# workflow-check-uncapped-loop

capability · `01M7DW5A3T5VC7PKSB03JP67BA` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

go test -trimpath ./internal/workflow -run TestPlainRedLoops covers a command step whose on.fail names the agent step before it (one loop), the same with handback:<agent step> (none), and a red route to a later step (none); go test -trimpath ./cmd/rhiza -run TestScripts/workflow-check-loop shows the warning in text and --json with exit 0, no such warning for ship-change, and none for a workflow from rhiza workflow new; go test -trimpath ./cmd/rhiza -run TestScripts/workflow-new passes. Touches internal/workflow/loops.go, internal/workflow/loops_test.go, internal/api/workflow_check.go, internal/api/workflow_new.go and cmd/rhiza/testdata/script/workflow-check-loop.txtar.

## statement

rhiza workflow check warns for each plain red loop: a step's on.<outcome> other than pass naming by plain step name a step that leads back to it through an agent step, naming the step, the outcome and its target and suggesting handback:<agent step>, which handback_cap counts; the workflow new skeleton routes its test failure with handback:build, and the warning never changes valid or the exit code.

## Links

### Out

- decided-by [uncapped-loop-is-plain-red-loop](../../../../decision/uncapped-loop-is-plain-red-loop.md)
- decided-by [skill-check-plan-shape](../../../../decision/skill-check-plan-shape.md)

### In

- verifies [workflow-check-uncapped-loop-verified](../../../../evidence/workflow-check-uncapped-loop-verified.md)
- verifies [workflow-check-uncapped-loop-qa-1](../../../../evidence/workflow-check-uncapped-loop-qa-1.md)
- verifies [workflow-check-uncapped-loop-qa-2](../../../../evidence/workflow-check-uncapped-loop-qa-2.md)
- verifies [workflow-check-uncapped-loop-qa-3](../../../../evidence/workflow-check-uncapped-loop-qa-3.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
