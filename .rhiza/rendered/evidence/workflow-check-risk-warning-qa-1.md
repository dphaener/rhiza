# workflow-check-risk-warning-qa-1

evidence · `01M6W9QDWVJPX4SRFKB7Q2FXZ8` · actor [claude-code-afe51521](../actor/claude-code-afe51521.md)

## result

PASS

## statement

QA at bbfe87e3 with a built binary (go build -trimpath) in a scratch Rhiza root and project, real git config, no fakes: rhiza workflow check plan with no bypass: 'workflow plan is valid (builtin:plan)', no warning. With plan auto + bypassPermissions: exit 0, valid line plus '  warning: workflow plan runs from a queued request with bypassed permissions in interview, review, critic, revise, arbiter: outside text reaches agents with no permission checks'; --json valid true, problems [], warnings holds that line. ship-change gets its warning only while plan, accept and ship are all auto; acceptEdits clears it.

## Links

### Out

- verifies [workflow-check-risk-warning](../capability/rhiza/workflows/workflow-authoring/workflow-check-risk-warning.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
