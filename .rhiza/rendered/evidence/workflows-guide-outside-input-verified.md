# workflows-guide-outside-input-verified

evidence · `01M6W7DSNGXYW7WNX7CWK2JZV0` · actor [claude-code-6c9df63d](../actor/claude-code-6c9df63d.md)

## result

PASS

## statement

At 09b17f17 ran go test -trimpath -count=1 ./internal/guide/...: ok (TestWorkflowsText's 'outside input' case checks ## Outside input, bypassPermissions, rhiza.autonomy.plan, rhiza request, 'Rhiza cannot see', permission_mode, rhiza.autonomy.signoff, plus 'does not gate intake', rhiza.daemon.permissionMode, 'outside Rhiza', 'accept the plan', "plan workflow's steps"; the guide's existing size checks pass in the same run). Built the commit and read `rhiza guide workflows`: the Outside input section says the warning means a request's text and --ref reach an autonomous bypassPermissions workflow, every queued request counts, Rhiza cannot see what a scheduled rhiza request sends, nothing is refused; per decision outside-input-plan-owner-not-intake it states plan owner does not gate intake (planned at once under rhiza.daemon.permissionMode, owner only accepts the plan) and lists the mitigations: intake gated by people outside Rhiza, narrower rhiza.daemon.permissionMode or permission_mode on the plan workflow's steps, plan owner and signoff owner for what runs after. This fixes attempt 1's finding f1. The fix commit stages only its builder's own actor record.

## Links

### Out

- verifies [workflows-guide-outside-input](../capability/rhiza/workflows/workflow-authoring/workflows-guide-outside-input.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
