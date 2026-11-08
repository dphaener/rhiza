# workflows-guide-outside-input

capability · `01M6W4ZN88P5XT2RXYH35YRVMZ` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

A test in internal/guide (go test -trimpath ./internal/guide) checks the workflows topic has the Outside input section naming bypassPermissions, rhiza.autonomy.plan (as not gating intake), rhiza.daemon.permissionMode, rhiza request, the schedule caveat, permission_mode and rhiza.autonomy.signoff, and the guide's existing size checks still pass

## statement

The workflows guide topic has an Outside input section: what the risk warning means (a request's text and --ref reach autonomous agents running in a bypass mode), that every queued request counts and Rhiza cannot see what a scheduled command running rhiza request does, that rhiza.autonomy.plan owner does not gate intake (each request is planned at once under rhiza.daemon.permissionMode; the owner only accepts the plan before the rest runs), and the usual mitigations: people outside Rhiza gate intake (never a schedule), a narrower rhiza.daemon.permissionMode or permission_mode on the plan workflow's steps, and rhiza.autonomy.plan owner and rhiza.autonomy.signoff owner for what runs after; nothing is refused

## Links

### Out

- decided-by [risk-warnings-shape](../../../../decision/risk-warnings-shape.md)

### In

- blocks [review-01m6w4k1-r3-f2](../../../../question/review-01m6w4k1-r3-f2.md)
- verifies [workflows-guide-outside-input-verified](../../../../evidence/workflows-guide-outside-input-verified.md)
- verifies [workflows-guide-outside-input-qa-1](../../../../evidence/workflows-guide-outside-input-qa-1.md)
- verifies [workflows-guide-outside-input-verified-2](../../../../evidence/workflows-guide-outside-input-verified-2.md)
- verifies [pr-365-risk-warnings](../../../../evidence/pr-365-risk-warnings.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
