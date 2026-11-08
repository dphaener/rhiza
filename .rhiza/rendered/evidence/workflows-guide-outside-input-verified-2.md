# workflows-guide-outside-input-verified-2

evidence · `01M6WA3M70RBGMPEZ7MJZYJ2ED` · actor [claude-code-986b5c73](../actor/claude-code-986b5c73.md)

## result

PASS

## statement

At 43079b55 ran go test -trimpath -count=1 ./internal/guide/...: ok. TestWorkflowsText passes and checks the workflows topic's '## Outside input' section for bypassPermissions, rhiza.autonomy.plan with 'does not gate intake', rhiza.daemon.permissionMode, rhiza request, the schedule caveat ('Rhiza cannot see' what a scheduled rhiza request sends), permission_mode on the plan workflow's steps, and rhiza.autonomy.signoff; TestCaps and TestCapCatchesOversize (size checks) pass. Read internal/guide/text/workflows.md: the section states the warning, every queued request counts, nothing is refused, plan owner only accepts the plan, and the mitigations (people outside Rhiza gate intake, never a schedule; narrower permission mode; plan/signoff owner).

## Links

### Out

- verifies [workflows-guide-outside-input](../capability/rhiza/workflows/workflow-authoring/workflows-guide-outside-input.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
