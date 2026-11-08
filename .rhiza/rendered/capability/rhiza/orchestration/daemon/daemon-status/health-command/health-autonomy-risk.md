# health-autonomy-risk

capability · `01M6W4ZN34KQZG9YEWQ1TKCXM3` · parent [health-command](../health-command.md)

## acceptance

Tests in internal/api for the autonomy part cover a flagged plan (warn, summary phrase, one detail line naming plan, request and its steps, Next rhiza guide workflows), nothing flagged (ok, unchanged), and flagged while paused (warn, paused Next kept); go test -trimpath ./cmd/rhiza -run TestScripts/health-risk-warning shows rhiza health exiting 0 with the autonomy part warn

## statement

When any workflow is on the request-exposed list, rhiza health's autonomy part is warn (never bad), its summary gains 'outside text reaches agents with no permission checks', it has one detail line per flagged workflow naming it, its trigger and its bypassing steps, and its Next is rhiza guide workflows unless the daemon is paused, which keeps the paused Next

## Links

### Out

- depends-on [request-exposed-workflows](../../../../workflows/autonomous-runs/request-exposed-workflows.md)
- decided-by [risk-warnings-shape](../../../../../../decision/risk-warnings-shape.md)

### In

- verifies [health-autonomy-risk-verified](../../../../../../evidence/health-autonomy-risk-verified.md)
- verifies [health-autonomy-risk-qa-1](../../../../../../evidence/health-autonomy-risk-qa-1.md)
- verifies [pr-365-risk-warnings](../../../../../../evidence/pr-365-risk-warnings.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
