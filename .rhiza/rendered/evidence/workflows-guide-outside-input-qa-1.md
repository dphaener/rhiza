# workflows-guide-outside-input-qa-1

evidence · `01M6W9QE7MX2YM8TJQJHXMQ7T3` · actor [claude-code-afe51521](../actor/claude-code-afe51521.md)

## result

PASS

## statement

QA at bbfe87e3 with a built binary (go build -trimpath) in a scratch Rhiza root and project, real git config, no fakes: rhiza guide workflows prints the '## Outside input' section: warns on request text and --ref reaching an autonomous bypassPermissions workflow, every queued request counts, Rhiza cannot see what a scheduled rhiza request sends, nothing refused; says rhiza.autonomy.plan owner does not gate intake; mitigations: people not schedules queue requests, narrower rhiza.daemon.permissionMode or step permission_mode, plan owner and rhiza.autonomy.signoff owner for what runs after.

## Links

### Out

- verifies [workflows-guide-outside-input](../capability/rhiza/workflows/workflow-authoring/workflows-guide-outside-input.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
