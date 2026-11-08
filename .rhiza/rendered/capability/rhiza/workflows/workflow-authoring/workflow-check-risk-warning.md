# workflow-check-risk-warning

capability · `01M6W4ZMYKAC3V8D3C5E539YHV` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/workflow-check-risk passes: with git config rhiza.autonomy.plan auto and rhiza.daemon.permissionMode bypassPermissions, workflow check plan exits 0 with valid true and a warning (JSON warnings field and human line) naming plan, the request trigger and its steps; with rhiza.autonomy.plan owner it prints no warning; describe-golden is regenerated for the new response field

## statement

rhiza workflow check <name> reports, beside its problems, a warning when the named workflow is one the request-exposed list flags, naming the trigger and the bypassing steps and saying outside text reaches agents with no permission checks; valid and the exit code are unchanged by it

## Links

### Out

- depends-on [request-exposed-workflows](../autonomous-runs/request-exposed-workflows.md)
- decided-by [risk-warnings-shape](../../../../decision/risk-warnings-shape.md)

### In

- verifies [workflow-check-risk-warning-verified](../../../../evidence/workflow-check-risk-warning-verified.md)
- verifies [workflow-check-risk-warning-qa-1](../../../../evidence/workflow-check-risk-warning-qa-1.md)
- verifies [pr-365-risk-warnings](../../../../evidence/pr-365-risk-warnings.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
