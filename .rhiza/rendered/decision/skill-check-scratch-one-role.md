# skill-check-scratch-one-role

decision · `01M7DW59N7VRS3Z9CKN703EPCJ`

## question

Change set 01M79KHB: how should rhiza skill check run a skill against a fixture change set?

## reasoning

The intent says a skill passes rhiza skill check before it runs for real, so the check must not touch the user's repo.

## statement

rhiza skill check <invoke> --role builder|verifier|check [--agent] builds a scratch repo outside every git repository with a one-leaf fixture change set, runs the invoke as that role's session through a one-step workflow in the foreground, and reports each contract gate as pass or fail, exiting 1 on any fail, with --json giving the same report.

## wording

Scratch repo, one role (assumed, would have escalated)

## Links

### Out

- derived-from [plan-skill-check](../source/plan-skill-check.md)

### In

- decided-by [skill-check-builder](../capability/rhiza/workflows/workflow-authoring/skill-check-builder.md)
- decided-by [skill-check-verifier-check](../capability/rhiza/workflows/workflow-authoring/skill-check-verifier-check.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
