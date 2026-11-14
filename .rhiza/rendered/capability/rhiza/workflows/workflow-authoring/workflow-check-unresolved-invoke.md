# workflow-check-unresolved-invoke

capability · `01M7DW5A21WH33SFE7PWYMAJMV` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

go test -trimpath ./internal/api -run TestWorkflowCheckUnresolvedInvoke covers a builder role and an agent step whose skill is missing (one warning each, naming the step, role and command, valid stays true) and the same workflow with the skills installed (no warning); go test -trimpath ./cmd/rhiza -run TestScripts/workflow-check-invoke shows the warning in text and --json and exit 0. Touches internal/api/workflow_check.go, internal/api/session.go (sharing checkStartSkill's lookup), internal/api/workflow_check_invoke_test.go and cmd/rhiza/testdata/script/workflow-check-invoke.txtar.

## statement

rhiza workflow check warns, for each agent step and role whose invoke starts with a slash command that its agent's adapter finds installed neither in the agent's config folder nor in the repo, naming the step, the role and the command, and for a Rhiza skill the rhiza guide install --agent <agent> command; it never changes valid or the exit code.

## Links

### Out

- decided-by [skill-check-plan-shape](../../../../decision/skill-check-plan-shape.md)

### In

- verifies [workflow-check-unresolved-invoke-verified](../../../../evidence/workflow-check-unresolved-invoke-verified.md)
- verifies [workflow-check-unresolved-invoke-qa-1](../../../../evidence/workflow-check-unresolved-invoke-qa-1.md)
- verifies [workflow-check-unresolved-invoke-qa-2](../../../../evidence/workflow-check-unresolved-invoke-qa-2.md)
- verifies [workflow-check-unresolved-invoke-qa-3](../../../../evidence/workflow-check-unresolved-invoke-qa-3.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
