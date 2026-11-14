# workflow-check-merge-no-gate

capability · `01M7DW5A714DZXJS2M54B08WPY` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

go test -trimpath ./internal/workflow -run TestMergeWithoutGate covers a merge reached only through a gate step (no warning), a merge reachable around the gate (warning naming the step and the path) and a workflow with no merge (none); go test -trimpath ./cmd/rhiza -run TestScripts/workflow-check-merge shows the warning in text and --json with exit 0 and no such warning for ship-change. Touches internal/workflow/merge_gate.go, internal/workflow/merge_gate_test.go, internal/api/workflow_check.go and cmd/rhiza/testdata/script/workflow-check-merge.txtar.

## statement

rhiza workflow check warns when a merge step (a command step whose run calls rhiza ship merge or gh pr merge) can be reached from the start step along a path that enters no gate step, naming the merge step and one such path; it never changes valid or the exit code.

## Links

### Out

- decided-by [skill-check-plan-shape](../../../../decision/skill-check-plan-shape.md)

### In

- verifies [workflow-check-merge-no-gate-verified](../../../../evidence/workflow-check-merge-no-gate-verified.md)
- verifies [workflow-check-merge-no-gate-qa-1](../../../../evidence/workflow-check-merge-no-gate-qa-1.md)
- verifies [workflow-check-merge-no-gate-qa-2](../../../../evidence/workflow-check-merge-no-gate-qa-2.md)
- verifies [workflow-check-merge-no-gate-qa-3](../../../../evidence/workflow-check-merge-no-gate-qa-3.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
