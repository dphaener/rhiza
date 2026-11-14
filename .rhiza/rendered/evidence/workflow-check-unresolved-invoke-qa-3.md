# workflow-check-unresolved-invoke-qa-3

evidence · `01M7E42SDK4FQ82W340KQAX2H3` · actor [claude-code-2248efc5](../actor/claude-code-2248efc5.md)

## result

PASS

## statement

QA cycle 3: with ~/.claude removed, rhiza workflow check ship-change warns that step conduct's /rhiza-conduct is not installed with rhiza guide install --agent claude, in text and --json, exit 0; with the skills reinstalled, no warning. TestScripts workflow-check-invoke passes.

## Links

### Out

- verifies [workflow-check-unresolved-invoke](../capability/rhiza/workflows/workflow-authoring/workflow-check-unresolved-invoke.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
