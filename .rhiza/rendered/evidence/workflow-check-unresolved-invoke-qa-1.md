# workflow-check-unresolved-invoke-qa-1

evidence · `01M7E24V0HQPVFKFWS6DRT1GGD` · actor [claude-code-98071bd0](../actor/claude-code-98071bd0.md)

## result

PASS

## statement

QA: rhiza workflow check on a workflow with an agent step /nosuch-skill and a leaves builder /rhiza-nosuchbuilder warns for each, naming step, role and command, exit 0, valid true; with HOME skills removed, the /rhiza-verifier role also warns and adds run rhiza guide install --agent claude; a repo-local .claude/skills/nosuch-skill clears its warning; installed skills give none.

## Links

### Out

- verifies [workflow-check-unresolved-invoke](../capability/rhiza/workflows/workflow-authoring/workflow-check-unresolved-invoke.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
